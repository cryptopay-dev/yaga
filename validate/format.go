package validate

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/cryptopay-dev/yaga/errors"
	"gopkg.in/go-playground/validator.v9"
)

type tagParser = func(tag reflect.StructTag) string

// fieldName parse struct for field name in json / form / query tags:
func fieldName(v reflect.Value, field string) string {
	var (
		tp      = v.Type()
		options = []tagParser{
			// Parse json-tag
			func(tag reflect.StructTag) string {
				val := tag.Get("json")
				return strings.Split(val, ",")[0]
			},

			// Parse form-tag
			func(tag reflect.StructTag) string {
				val := tag.Get("form")
				return strings.Split(val, ",")[0]
			},

			// Parse query-tag
			func(tag reflect.StructTag) string {
				val := tag.Get("query")
				return strings.Split(val, ",")[0]
			},
		}
	)

	if f, ok := tp.FieldByName(field); ok {
		for _, o := range options {
			if val := o(f.Tag); len(val) > 0 {
				return val
			}
		}
	}

	return strings.ToLower(field)
}

// Options to call CheckErrors method
type Options struct {
	Struct    interface{}
	Errors    error
	Formatter func(fields []string) string
}

// defaultFormatter generates "bad `field1`, `field2`"
func defaultFormatter(fields []string) string {
	return "bad `" + strings.Join(fields, "`, `") + "`"
}

// CheckErrors of validator and return formatted errors:
func CheckErrors(opts Options) (ok bool, err error) {
	var fieldsErr validator.ValidationErrors

	if opts.Struct == nil || opts.Errors == nil {
		return
	}

	if opts.Formatter == nil {
		opts.Formatter = defaultFormatter
	}

	if fieldsErr, ok = opts.Errors.(validator.ValidationErrors); ok {
		var (
			fields = make([]string, 0, len(fieldsErr))
			val    = reflect.ValueOf(opts.Struct)
		)

		if val.Kind() == reflect.Ptr && !val.IsNil() {
			val = val.Elem()
		}

		for _, field := range fieldsErr {
			fields = append(fields, fieldName(val, field.Field()))
		}

		err = errors.NewError(http.StatusBadRequest, opts.Formatter(fields))
	}

	return
}
