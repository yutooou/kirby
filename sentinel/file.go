package sentinel

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/yutooou/kirby/models"
	"path/filepath"
	"sync"
)

type FileSystemSentinel struct {
	path               string
	allKirbyFile       map[string]KirbyFile
	acceptFileProtocol []string
}

func NewFileSystemSentinel(path string) FileSystemSentinel {
	return FileSystemSentinel{
		path:               path,
		allKirbyFile:       make(map[string]KirbyFile),
		acceptFileProtocol: []string{".json", ".kbl"},
	}
}

type KirbyFile struct {
	path  string
	name  string
	kirby models.Kirby
	md5   string
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

func (f FileSystemSentinel) inAcceptFileProtocol(path string) bool {
	ext := filepath.Ext(path)
	for _, v := range f.acceptFileProtocol {
		if v == ext {
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
