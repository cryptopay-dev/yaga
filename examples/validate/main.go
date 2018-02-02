package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/validate"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

// Request struct for example
type Request struct {
	A string `form:"a" json:"a" validate:"required"`
	B int64  `form:"a" json:"a" validate:"required,gt=1"`
}

func main() {
	e := echo.New()
	v := validator.New()

	// ... if need connect custom validators
	// (see - https://godoc.org/gopkg.in/go-playground/validator.v9#hdr-Custom_Validation_Functions)

	// Creates echo-like validator
	e.Validator = validate.New(v)

	// Our test action
	e.POST("/", func(ctx echo.Context) error {
		var req Request

		// Try to parse request
		if err := ctx.Bind(&req); err != nil {
			return err
		}

		// Try to validate request
		if err := ctx.Validate(&req); err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, req)
	})

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
