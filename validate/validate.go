package validate

import (
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

// validate is wrapper of echo.Validator
type validate struct {
	Validator *validator.Validate
}

// New creates new wrapper of validator for echo.Validator
func New(v *validator.Validate) echo.Validator {
	return &validate{v}
}

// Validate a struct(s) exposed fields, and automatically validates nested struct(s), unless otherwise specified.
func (v *validate) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}
