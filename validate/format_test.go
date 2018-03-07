package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

type test1 struct {
	A int `json:"a_custom" validate:"gt=0"`
	B int `form:"b_custom" validate:"required"`
	C int
	D int `query:"someValue" validate:"required"`
	E int `json:"-" form:"-" query:"-" param:"-" xml:"-" yaml:"-" validate:"required"`
	F int `json:"-" param:"f_custom" xml:"-" yaml:"-" validate:"required"`
	G int `xml:"g_custom" yaml:"-" validate:"required"`
	H int `yaml:"h_custom" validate:"required"`
}

var testCases = []struct {
	Struct interface{}
	Error  error
}{
	{
		Struct: test1{A: -1},
		Error:  newError(400, "bad `a_custom`, `b_custom`, `someValue`, `e`, `f_custom`, `g_custom`, `h_custom`"),
	},

	{
		Struct: test1{A: 1, B: 1, D: 1, E: 1, F: 1, G: 1, H: 1},
		Error:  nil,
	},
}

func TestCheckErrors(t *testing.T) {

	v := validator.New()

	for _, test := range testCases {
		errValidate := v.Struct(test.Struct)

		ok, err := CheckErrors(Options{
			Struct: test.Struct,
			Errors: errValidate,
		})

		if test.Error != nil {
			assert.True(t, ok)
			assert.Error(t, err)
			assert.Equal(t, test.Error.Error(), err.Error())
		} else {
			assert.False(t, ok)
			assert.NoError(t, err)
		}
	}
}
