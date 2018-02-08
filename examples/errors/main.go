package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	log := zap.New(zap.Development)

	e := web.New(web.Options{
		Logger: log,
	})

	logic, err := errors.New(errors.Options{
		Debug:  true,
		Logger: log,
	})

	if err != nil {
		panic(err)
	}

	e.HTTPErrorHandler = logic.Capture
	e.Use(logic.Recover())

	e.GET("/", func(ctx web.Context) error {
		return ctx.String(http.StatusOK, "Hello world")
	})

	e.GET("/bad-request", func(ctx web.Context) error {
		return errors.NewError(http.StatusBadRequest, "Bad request")
	})

	e.GET("/formatted-error", func(ctx web.Context) error {
		return errors.NewErrorf(http.StatusBadRequest, "Bad request '%s'", "formatted")
	})

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
