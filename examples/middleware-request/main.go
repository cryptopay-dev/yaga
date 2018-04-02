package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	log.Init()

	e, err := web.New(web.Options{})
	if err != nil {
		log.Panic(err)
	}

	e.Use(web.RayTraceID)

	e.GET("/", func(ctx web.Context) error {
		return ctx.String(
			http.StatusOK,
			ctx.Request().Header.Get(web.RayTraceHeader),
		)
	})

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
