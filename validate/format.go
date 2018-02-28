package validate

import (
	"net/http"
	"strings"

	"github.com/cryptopay-dev/yaga/errors"
	"gopkg.in/go-playground/validator.v9"
)

// CheckErrors of validator and return formatted errors:
func CheckErrors(valErr error) (ok bool, err error) {
	var fieldsErr validator.ValidationErrors
	if fieldsErr, ok = valErr.(validator.ValidationErrors); ok {
		var fields []string
		for _, field := range fieldsErr {
			fields = append(fields, strings.ToLower(
				field.Field(),
			))
		}
		err = errors.NewError(http.StatusBadRequest, "bad "+strings.Join(fields, ", "))
	}
	return
}
