package errors

import "fmt"

type LogicError struct {
	Code    int
	Message string
}

func (err *LogicError) Error() string {
	return err.Message
}

func NewError(code int, message string) *LogicError {
	return &LogicError{code, message}
}

func NewErrorf(code int, format string, a ...interface{}) *LogicError {
	return &LogicError{code, fmt.Sprintf(format, a)}
}
