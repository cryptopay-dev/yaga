package errors

import "fmt"

// LogicError used for web-controllers
type LogicError struct {
	Code    int
	Message string
}

// Error is an implementation of error interface
func (err *LogicError) Error() string {
	return err.Message
}

// NewError return LogicError with http-code and error message
func NewError(code int, message string) *LogicError {
	return &LogicError{code, message}
}

// NewErrorf return LogicError with http-code and formatted error message
func NewErrorf(code int, format string, a ...interface{}) *LogicError {
	return &LogicError{code, fmt.Sprintf(format, a)}
}
