package web

import "fmt"

// Error used for web-controllers
type Error struct {
	Code    int
	Message string
}

// Error is an implementation of error interface
func (err *Error) Error() string {
	return err.Message
}

// NewError return Error with http-code and error message
func NewError(code int, message string) *Error {
	return &Error{code, message}
}

// NewErrorf return Error with http-code and formatted error message
func NewErrorf(code int, format string, a ...interface{}) *Error {
	return &Error{code, fmt.Sprintf(format, a...)}
}
