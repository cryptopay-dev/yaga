package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/cryptopay-dev/yaga/middlewares/request"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Logger = nop.New()
	e.HideBanner = true

	e.Use(request.RayTraceID(zap.New(zap.Development)))

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(
			http.StatusOK,
			ctx.Request().Header.Get(request.RayTraceHeader),
		)
	})

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
