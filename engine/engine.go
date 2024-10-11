package engine

import "github.com/yutooou/kirby/models"

type Engine interface {
	Run() (kch chan models.KirbyModel, ech chan error)
}
