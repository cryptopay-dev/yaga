package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Logger = nop.New()
	e.HideBanner = true

	logic, err := errors.New(errors.Options{
		Debug:  true,
		Logger: zap.New(zap.Development),
	})

	if err != nil {
		panic(err)
	}

	e.HTTPErrorHandler = logic.Capture
	e.Use(logic.Recover())

	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello world")
	})

	e.GET("/bad-request", func(ctx echo.Context) error {
		return errors.NewError(http.StatusBadRequest, "Bad request")
	})

	e.GET("/formatted-error", func(ctx echo.Context) error {
		return errors.NewErrorf(http.StatusBadRequest, "Bad request '%s'", "formatted")
	})

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
