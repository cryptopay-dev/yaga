package validate

import (
	"testing"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

func TestCheckErrors(t *testing.T) {
	var test struct {
		A int `validate:"gt=0"`
		B int `validate:"required"`
		C int
	}

	test.A = -1

	v := validator.New()

	errValidate := v.Struct(test)

	ok, err := CheckErrors(errValidate)

	assert.True(t, ok)
	assert.IsType(t, &errors.LogicError{}, err)
	assert.Equal(t, "bad `a`,`b`", err.Error())
}
