package validate

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cryptopay-dev/yaga/errors"
	"gopkg.in/go-playground/validator.v9"
)

// CheckErrors of validator and return formatted errors:
func CheckErrors(valErr error) (ok bool, err error) {
	var fieldsErr validator.ValidationErrors
	if fieldsErr, ok = valErr.(validator.ValidationErrors); ok {
		var (
			fields strings.Builder
			count  = len(fieldsErr) - 1
		)

		if _, wErr := fields.WriteString("bad "); wErr != nil {
			ok = false
			return
		}

		for i, field := range fieldsErr {
			if _, wErr := fmt.Fprintf(&fields, "`%s`", strings.ToLower(
				field.Field(),
			)); wErr != nil {
				ok = false
				return
			}

			if i < count {
				if wErr := fields.WriteByte(','); wErr != nil {
					ok = false
					return
				}
			}
		}
		err = errors.NewError(http.StatusBadRequest, fields.String())
	}

	return
}
