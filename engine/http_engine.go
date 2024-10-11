package engine

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yutooou/kirby/config"
	"github.com/yutooou/kirby/models"
	"net/http"
	"path"
)

const (
	rootPath = "/_kirby"
	infoPath = "/_info"
)

var LocalHttpEngine *HttpEngine

type HttpEngine struct {
	addr string
	r    *gin.Engine
	srv  *http.Server
	kch  chan models.KirbyModel
	ech  chan error
}

func init() {
	addr := config.LocalConf.Engine.Http.Addr
	r := initHandler()
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	kch := make(chan models.KirbyModel)
	ech := make(chan error)
	LocalHttpEngine = &HttpEngine{
		addr: addr,
		r:    r,
		srv:  srv,
		kch:  kch,
		ech:  ech,
	}
}

func (h *HttpEngine) Run() (kch chan models.KirbyModel, ech chan error) {
	go h.serve()
	go h.listen()
	return h.kch, h.ech
}

func (h *HttpEngine) serve() {
	err := h.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		h.ech <- err
	}
}

func (h *HttpEngine) listen() {
	for {
		select {
		case model := <-h.kch:
			h.restart(model)
		}
	}
}

func (h *HttpEngine) stop() {
	err := h.srv.Shutdown(context.Background())
	if err != nil {
		h.ech <- err
	}
}

func (h *HttpEngine) restart(model models.KirbyModel) {
	h.stop()
	r := initHandler()
	buildHandler(model, r)
	h.r = r
	h.srv = &http.Server{
		Addr:    h.addr,
		Handler: r,
	}
	go h.serve()
}

func buildHandler(model models.KirbyModel, r *gin.Engine) {
	for _, kirby := range model {
		r.GET(path.Join(kirby.Info.Code, infoPath), showInfo(kirby.Info))
	}
}

func showKirbyInfo() gin.HandlerFunc {
	info := models.Info{
		Name:       "Kirby",
		Code:       rootPath,
		Version:    "0.0.1",
		Desc:       "Kirby, data simulation engine.",
		EngineType: models.HTTP,
	}
	return showInfo(info)
}

func showInfo(info models.Info) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, info)
	}
}

func initHandler() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET(path.Join(rootPath, infoPath), showKirbyInfo())
	return r
}
