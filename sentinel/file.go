package sentinel

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yutooou/kirby/models"
	"github.com/yutooou/kirby/utils"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileSystemSentinel struct {
	path               string
	allKirbyFile       map[string]KirbyFile
	acceptFileProtocol []KirbyFileProtocol
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
		allKirbyFile:       make(map[string]KirbyFile),
		acceptFileProtocol: []KirbyFileProtocol{JsonProtocol, KblProtocol},
	}
	sentinel.initKirbyFile()

	return sentinel
}

func (f FileSystemSentinel) Watch() (kch chan models.KirbyModel, ech chan error) {
	kch = make(chan models.KirbyModel)
	ech = make(chan error)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ech <- fmt.Errorf("watch file system error[init]. err: %v", err)
		return
	}
	defer watcher.Close()
	err = watcher.Add(f.path)
	if err != nil {
		ech <- fmt.Errorf("watch file system error[watch]. err: %v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		var err error
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
					err = f.modify(event.Name)
					if err != nil {
						ech <- fmt.Errorf("watch file system error[modify]. err: %v", err)
						continue
					}
					kch <- f.kirbyModel()
				case fsnotify.Remove:
					err = f.delete(event.Name)
					if err != nil {
						ech <- fmt.Errorf("watch file system error[delete]. err: %v", err)
						continue
					}
					kch <- f.kirbyModel()
				default:
					continue
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					goto exit
				}
				ech <- fmt.Errorf("watch file system error[file]. err: %v", err)
			}
		}
	exit:
		wg.Done()
	}()
	wg.Wait()
	return
}

func (f FileSystemSentinel) initKirbyFile() {
	filepath.Walk(f.path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !f.inAcceptFileProtocol(path) {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		md5, err := utils.FileMD5F(file)
		if err != nil {
			return err
		}
		ext := filepath.Ext(path)

		// 获取path除了 ext 之外的文件名
		code := strings.TrimSuffix(info.Name(), ext)

		kf := KirbyFile{
			path:     path,
			fileName: info.Name(),
			code:     code,
			protocol: KirbyFileProtocol(ext),
			md5:      md5,
			kirby: models.Kirby{
				Info: models.Info{
					Code: code,
				},
			},
		}
		f.allKirbyFile[code] = kf
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

func (f FileSystemSentinel) modify(path string) error {
	// todo
	panic("todo")
}

func (f FileSystemSentinel) delete(path string) error {
	// todo
	panic("todo")
}
