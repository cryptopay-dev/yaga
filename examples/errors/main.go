package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	log := zap.New(zap.Development)

	e, err := web.New(web.Options{
		Logger: log,
	})
	if err != nil {
		log.Panic(err)
	}

	e.GET("/", func(ctx web.Context) error {
		return ctx.String(http.StatusOK, "Hello world")
	})

	e.GET("/bad-request", func(ctx web.Context) error {
		return web.NewError(http.StatusBadRequest, "Bad request")
	})

	e.GET("/formatted-error", func(ctx web.Context) error {
		return web.NewErrorf(http.StatusBadRequest, "Bad request '%s'", "formatted")
	})

	if err := e.Start(":8080"); err != nil {
		log.Panic(err)
	}
}
