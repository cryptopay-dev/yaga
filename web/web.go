package web

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cryptopay-dev/go-metrics"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	errStartServerTpl    = "error while starting server: %v"
	startServerOnPortTpl = "server starting on port '%s'"
	emptyBindEnv         = "BIND env for server not set, use default port: '%s'"

	defaultBind = ":8080"
)

type recoverer interface {
	Capture(error, echo.Context)
	Recover() echo.MiddlewareFunc
}

// Options contains a parameters for new Echo instance.
type Options struct {
	Logger    echo.Logger
	Error     recoverer
	Debug     bool
	Validator echo.Validator
}

// Errors
var (
	ErrUnsupportedMediaType        = echo.ErrUnsupportedMediaType
	ErrNotFound                    = echo.ErrNotFound
	ErrUnauthorized                = echo.ErrUnauthorized
	ErrForbidden                   = echo.ErrForbidden
	ErrMethodNotAllowed            = echo.ErrMethodNotAllowed
	ErrStatusRequestEntityTooLarge = echo.ErrStatusRequestEntityTooLarge
	ErrValidatorNotRegistered      = echo.ErrValidatorNotRegistered
	ErrRendererNotRegistered       = echo.ErrRendererNotRegistered
	ErrInvalidRedirectCode         = echo.ErrInvalidRedirectCode
	ErrCookieNotFound              = echo.ErrCookieNotFound
)

var (
	// NewHTTPError creates a new HTTPError instance.
	NewHTTPError = echo.NewHTTPError

	initMetricsOnce = sync.Once{}
)

type (
	// Context from echo.Context (for shadowing)
	Context = echo.Context

	// HandlerFunc from echo.HandlerFunc (for shadowing)
	HandlerFunc = echo.HandlerFunc

	// MiddlewareFunc from echo.MiddlewareFunc (for shadowing)
	MiddlewareFunc = echo.MiddlewareFunc

	// Engine from echo.Echo (for shadowing)
	Engine = echo.Echo

	// Group from echo.Group (for shadowing)
	Group = echo.Group

	// BasicAuthValidator defines a function to validate BasicAuth credentials.
	BasicAuthValidator = middleware.BasicAuthValidator
)

// New creates an instance of Echo.
func New(opts Options) *Engine {
	e := echo.New()

	if opts.Logger != nil {
		e.Logger = opts.Logger
	}

	// TODO may be move to function?
	initMetricsOnce.Do(func() {
		// enable metrics
		if err := metrics.Setup(os.Getenv("METRICS_URL"), os.Getenv("METRICS_APP"), os.Getenv("METRICS_HOSTNAME")); err == nil {
			go func() {
				if errWatch := metrics.Watch(time.Second * 10); errWatch != nil {
					e.Logger.Errorf("Can't start watching for metrics: %v", errWatch)
				}
			}()
		} else {
			e.Logger.Error(err)
		}
	})

	e.Debug = opts.Debug
	e.HideBanner = true

	if opts.Validator != nil {
		e.Validator = opts.Validator
	}

	e.Binder = new(DefaultBinder)

	if opts.Error != nil {
		e.HTTPErrorHandler = opts.Error.Capture
		e.Use(opts.Error.Recover())
	}

	return e
}

// StartServer HTTP with custom address.
func StartServer(e *Engine, bind string) error {
	// start server
	if len(bind) == 0 {
		bind = defaultBind
		e.Logger.Infof(emptyBindEnv, bind)
	}

	e.Logger.Infof(startServerOnPortTpl, bind)
	if err := e.Start(bind); err != nil {
		return fmt.Errorf(errStartServerTpl, err)
	}

	return nil
}

// AddTrailingSlash returns a root level (before router) middleware which adds a
// trailing slash to the request `URL#Path`.
//
// Usage `Engine#Pre(AddTrailingSlash())`
func AddTrailingSlash() MiddlewareFunc {
	return middleware.AddTrailingSlashWithConfig(middleware.DefaultTrailingSlashConfig)
}

// RemoveTrailingSlash returns a root level (before router) middleware which removes
// a trailing slash from the request URI.
//
// Usage `Engine#Pre(RemoveTrailingSlash())`
func RemoveTrailingSlash() MiddlewareFunc {
	return middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{})
}

// BasicAuth returns an BasicAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func BasicAuth(fn BasicAuthValidator) MiddlewareFunc {
	c := middleware.DefaultBasicAuthConfig
	c.Validator = fn
	return middleware.BasicAuthWithConfig(c)
}
