package sentinel

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yutooou/kirby/models"
	"github.com/yutooou/kirby/utils"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileSystemSentinel struct {
	path               string
	allKirbyFile       map[string]*KirbyFile
	acceptFileProtocol []KirbyFileProtocol
	kch                chan models.KirbyModel
	ech                chan error
}

type KirbyFile struct {
	path     string
	fileName string
	code     string
	protocol KirbyFileProtocol
	kirby    models.Kirby
	md5      string
}

type KirbyFileProtocol string

const (
	JsonProtocol KirbyFileProtocol = ".json"
	KblProtocol  KirbyFileProtocol = ".kbl"
)

func NewFileSystemSentinel(path string) FileSystemSentinel {
	if path == "" {
		panic("path is empty, can't init file system sentinel")
	}
	sentinel := FileSystemSentinel{
		path:               path,
		allKirbyFile:       make(map[string]*KirbyFile),
		acceptFileProtocol: []KirbyFileProtocol{JsonProtocol, KblProtocol},
		kch:                make(chan models.KirbyModel),
		ech:                make(chan error, 1),
	}
	sentinel.initAllKirbyFile()

	return sentinel
}

func NewKirbyFile(path string) (*KirbyFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	md5, err := utils.FileMD5F(file)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(path)
	fileName := filepath.Base(path)

	// 获取path除了 ext 之外的文件名
	code := strings.TrimSuffix(fileName, ext)

	kf := &KirbyFile{
		path:     path,
		fileName: fileName,
		code:     code,
		protocol: KirbyFileProtocol(ext),
		md5:      md5,
		kirby: models.Kirby{
			Info: models.Info{
				Code: code,
			},
		},
	}
	return kf, nil
}

func (f FileSystemSentinel) Watch() (chan models.KirbyModel, chan error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		f.ech <- fmt.Errorf("watch file system error[init]. err: %v", err)
		return f.kch, f.ech
	}
	defer watcher.Close()
	err = watcher.Add(f.path)
	if err != nil {
		f.ech <- fmt.Errorf("watch file system error[watch]. err: %v", err)
		return f.kch, f.ech
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					goto exit
				}
				// 判断  filepath.Ext(event.Name) 文件后缀是不是在 acceptFileProtocol 中
				if !f.inAcceptFileProtocol(event.Name) {
					continue
				}
				switch event.Op {
				case fsnotify.Create | fsnotify.Write:
					if ok, err := f.modify(event.Name); ok {
						f.kch <- f.kirbyModel()
					} else if err != nil {
						f.ech <- fmt.Errorf("watch file system error[modify]. err: %v", err)
						continue
					}
				case fsnotify.Remove:
					if ok, err := f.delete(event.Name); ok {
						f.kch <- f.kirbyModel()
					} else if err != nil {
						f.ech <- fmt.Errorf("watch file system error[delete]. err: %v", err)
						continue
					}
				default:
					continue
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					goto exit
				}
				f.ech <- fmt.Errorf("watch file system error[file]. err: %v", err)
			}
		}
	exit:
		wg.Done()
	}()
	wg.Wait()
	return f.kch, f.ech
}

func (f FileSystemSentinel) initAllKirbyFile() {
	filepath.Walk(f.path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !f.inAcceptFileProtocol(path) {
			return nil
		}
		kf, err := NewKirbyFile(path)
		if err != nil {
			return err
		}
		f.allKirbyFile[kf.code] = kf
		return nil
	})
}

func (f FileSystemSentinel) inAcceptFileProtocol(path string) bool {
	ext := filepath.Ext(path)
	for _, v := range f.acceptFileProtocol {
		if v == KirbyFileProtocol(ext) {
			return true
		}
	}
	return false
}

func (f FileSystemSentinel) kirbyModel() models.KirbyModel {
	ret := make(models.KirbyModel, len(f.allKirbyFile))
	for _, kirbyFile := range f.allKirbyFile {
		ret = append(ret, kirbyFile.kirby)
	}
	return ret
}

func (f FileSystemSentinel) modify(path string) (ok bool, err error) {
	kf, err := NewKirbyFile(path)
	if err != nil {
		return false, fmt.Errorf("modify file err: %v", err)
	}
	if f.allKirbyFile[kf.code] != nil && f.allKirbyFile[kf.code].md5 == kf.md5 {
		log.Println("file not change. code: " + kf.code)
		return false, nil
	}
	f.allKirbyFile[kf.code] = kf
	return true, nil
}

func (f FileSystemSentinel) delete(path string) (ok bool, err error) {
	ext := filepath.Ext(path)
	fileName := filepath.Base(path)
	code := strings.TrimSuffix(fileName, ext)
	if _, ok := f.allKirbyFile[code]; !ok {
		return false, nil
	}
	delete(f.allKirbyFile, code)
	return true, nil
}
