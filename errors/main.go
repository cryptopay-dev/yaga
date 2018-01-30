package errors

import (
	"errors"
	"net/http"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/cryptopay-dev/yaga/middlewares/request"
	"github.com/cryptopay-dev/yaga/tracer"
	"github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

// Options for creating new Logic
type Options struct {
	Debug  bool
	Logger logger.Logger
}

// Response answer structure
type Response struct {
	Error  string   `json:"error"`
	Stack  []string `json:"stack"`
	Result []string `json:"result"`
}

// Logic is an structure for capture and recover
// web-errors and panics
type Logic struct {
	Opts Options
}

var (
	// ErrorEmptyLogger issued when logger received empty logger to New-method
	ErrorEmptyLogger = errors.New("options hasn't logger")
)

// New creates instance of Logic structure
func New(opts Options) (*Logic, error) {
	if opts.Logger == nil {
		return nil, ErrorEmptyLogger
	}

	e := new(Logic)
	e.Opts = opts

	return e, nil
}

// Capture web-errors and formatting answer
func (c *Logic) Capture(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	result := make([]string, 0)
	trace := make([]string, 0)
	message := err.Error()

	switch custom := err.(type) {
	case *echo.HTTPError:
		code = custom.Code
		message = custom.Message.(string)
	case *LogicError:
		code = custom.Code
	default:
		message = http.StatusText(code)
	}

	// Capture errors:
	if code >= http.StatusInternalServerError {
		raven.CaptureErrorAndWait(err, request.TraceTag(ctx))
		c.Opts.Logger.Error("Request error", zap.Error(err), request.TraceField(ctx))
	}

	// Capture stack trace:
	if c.Opts.Debug {
		trace = append(trace, tracer.Stack(err)...)
	}

	if errJSON := ctx.JSON(code, Response{
		Error:  message,
		Stack:  trace,
		Result: result,
	}); errJSON != nil {
		c.Opts.Logger.Errorf("LogicError.Capture error: %v", err)
	}
}
