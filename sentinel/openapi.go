package sentinel

import (
	"github.com/yutooou/kirby/models"
)

type OpenapiSentinel struct {
	addr map[string]string
}

func (f OpenapiSentinel) Watch() (kch chan models.KirbyModel, ech chan error) {
	kch = make(chan models.KirbyModel)
	ech = make(chan error)
	return
}
