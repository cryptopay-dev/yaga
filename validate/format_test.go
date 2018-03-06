package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

func TestCheckErrors(t *testing.T) {
	var test struct {
		A int `json:"a_custom" validate:"gt=0"`
		B int `form:"b_custom" validate:"required"`
		C int
		D int `query:"someValue" validate:"required"`
		E int `validate:"required"`
	}

	test.A = -1

	v := validator.New()

	errValidate := v.Struct(test)

	ok, err := CheckErrors(Options{
		Struct: test,
		Errors: errValidate,
	})

	assert.True(t, ok)
	assert.IsType(t, Error{}, err)
	assert.Equal(t, "bad `a_custom`, `b_custom`, `someValue`, `e`", err.Error())
}
