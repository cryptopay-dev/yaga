package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/validate"
	"github.com/cryptopay-dev/yaga/web"
	"gopkg.in/go-playground/validator.v9"
)

// Request struct for example
type Request struct {
	A string `form:"a" json:"a" validate:"required"`
	B int64  `form:"a" json:"a" validate:"required,gt=1"`
}

func main() {
	v := validator.New()

	// ... if need connect custom validators
	// (see - https://godoc.org/gopkg.in/go-playground/validator.v9#hdr-Custom_Validation_Functions)

	e := web.New(web.Options{
		// Creates echo-like validator
		Validator: validate.New(v),
	})

	// Our test action
	e.POST("/", func(ctx web.Context) error {
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
