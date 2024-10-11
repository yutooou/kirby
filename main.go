package main

import (
	"github.com/yutooou/kirby/engine"
)

func main() {
	_, httpEch := engine.LocalHttpEngine.Run()
	select {
	case e := <-httpEch:
		panic(e)
	}
}
