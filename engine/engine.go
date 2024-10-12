package engine

import "github.com/yutooou/kirby/models"

// Engine interface
// Engine is the interface that Kirby Engine should implement.
// such as http engine. rpc engine. syslog engine.
type Engine interface {
	Run() (kch chan models.KirbyModel, ech chan error)
}
