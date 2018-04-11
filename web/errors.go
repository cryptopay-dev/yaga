package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"sync"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/cryptopay-dev/yaga/tracer"
	"github.com/cryptopay-dev/yaga/validate"
	"github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// LogicResponse answer structure
type LogicResponse struct {
	Error  string   `json:"error"`
	Stack  []string `json:"stack"`
	Result []string `json:"result"`
}

// Logic is an structure for capture and recover
// web-errors and panics
type Logic struct {
	Debug bool
}

var initRavenOnce = sync.Once{}

// NewLogic creates instance of Logic structure
func NewLogic(debug bool) *Logic {
	e := new(Logic)
	e.Debug = debug

	initRavenOnce.Do(func() {
		if err := raven.SetDSN(config.GetString("sentry_dsn")); err != nil {
			log.Error(err)
		}
	})

	return e
}

// CustomError interface
type CustomError interface {
	FormatResponse(ctx echo.Context)
}

// Capture web-errors and formatting answer
func (c *Logic) Capture(err error, ctx echo.Context) {
	code := http.StatusBadRequest
	result := make([]string, 0)
	trace := make([]string, 0)
	message := err.Error()

	switch custom := err.(type) {
	case CustomError:
		custom.FormatResponse(ctx)
		return
	case *json.UnmarshalTypeError:
		message = fmt.Sprintf("JSON parse error: expected=%v, got=%v, offset=%v", custom.Type, custom.Value, custom.Offset)
	case *json.SyntaxError:
		message = fmt.Sprintf("JSON parse error: offset=%v, error=%v", custom.Offset, custom.Error())
	case *xml.UnsupportedTypeError:
		message = fmt.Sprintf("XML parse error: type=%v, error=%v", custom.Type, custom.Error())
	case *xml.SyntaxError:
		message = fmt.Sprintf("XML parse error: line=%v, error=%v", custom.Line, custom.Error())
	case *echo.HTTPError:
		code = custom.Code
		message = custom.Message.(string)
	case *Error:
		code = custom.Code
	case validate.Error:
		code = custom.Code
	default:
		message = http.StatusText(code)
	}

	// Capture errors:
	if code >= http.StatusInternalServerError {
		raven.CaptureErrorAndWait(err, TraceTag(ctx))
		log.Error("Request error", zap.Error(err), TraceTag(ctx))
	}

	// Capture stack trace:
	if c.Debug {
		trace = append(trace, tracer.Stack(err)...)
	}

	if errJSON := ctx.JSON(code, LogicResponse{
		Error:  message,
		Stack:  trace,
		Result: result,
	}); errJSON != nil {
		log.Errorf("Error.Capture error: %v", err)
	}
}
