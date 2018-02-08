package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/cryptopay-dev/yaga/middlewares/request"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	log := zap.New(zap.Development)

	e := web.New(web.Options{
		Logger: log,
	})

	e.Use(request.RayTraceID(log))

	e.GET("/", func(ctx web.Context) error {
		return ctx.String(
			http.StatusOK,
			ctx.Request().Header.Get(request.RayTraceHeader),
		)
	})

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
