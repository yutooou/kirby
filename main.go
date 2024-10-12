package main

import (
	"github.com/yutooou/kirby/engine"
	"github.com/yutooou/kirby/sentinel"
)

func main() {
	httpKch, httpEch := engine.LocalHttpEngine.Run()

	kch := sentinel.RunAllSentinel()
	go func() {
		for {
			select {
			case k := <-kch:
				httpKch <- k
			}
		}
	}()

	select {
	case e := <-httpEch:
		panic(e)
	}
}
