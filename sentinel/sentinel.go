package sentinel

import (
	"github.com/yutooou/kirby/config"
	"github.com/yutooou/kirby/models"
	"log"
)

var (
	AllSentinels = make(map[string]Sentinel)
)

func init() {
	if config.LocalConf.Sentinel.File.Enable {
		AllSentinels["file"] = NewFileSystemSentinel(config.LocalConf.Sentinel.File.Dir)
	}
}

// Sentinel Interfaces that need to be implemented to watch configuration changes to change engine behavior
type Sentinel interface {
	// Watch returns a channel that will receive messages from the sentinel
	// this method must be non blocking
	Watch() (kch chan models.KirbyModel, ech chan error)
}

// RunAllSentinel runs all sentinels
// receive messages from all sentinels and send out
func RunAllSentinel() (kch chan models.KirbyModel) {
	kch = make(chan models.KirbyModel)
	for _, sentinel := range AllSentinels {
		modelsChan, errorsChan := sentinel.Watch()
		go func() {
			for {
				select {
				case kirbyModel := <-modelsChan:
					kch <- kirbyModel
				case err := <-errorsChan:
					log.Printf("sentinel error: %v\n", err)
				}
			}
		}()
	}
	return kch
}
