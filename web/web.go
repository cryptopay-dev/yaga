package web

import (
	"fmt"
	"os"
	"time"

	"github.com/cryptopay-dev/go-metrics"
	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/getsentry/raven-go"
	"github.com/labstack/echo"
)

const (
	errStartServerTpl    = "error while starting server: %v"
	startServerOnPortTpl = "server starting on port '%s'"
	emptyBindEnv         = "BIND env for server not set, use default port: '%s'"

	defaultBind = ":8080"
)

// Options contains a parameters for new Echo instance.
type Options struct {
	// TODO suggest change to echo.Logger
	Logger    logger.Logger
	Error     errors.Logic
	Debug     bool
	Validator echo.Validator
}

// Context from echo.Context (for shadowing)
type Context = echo.Context

// New creates an instance of Echo.
func New(opts Options) *echo.Echo {
	if err := raven.SetDSN(os.Getenv("SENTRY_DSN")); err != nil {
		opts.Logger.Error(err)
	}

	// enable metrics
	if err := metrics.Setup(os.Getenv("METRICS_URL"), os.Getenv("METRICS_APP"), os.Getenv("METRICS_HOSTNAME")); err == nil {
		go func() {
			if errWatch := metrics.Watch(time.Second * 10); errWatch != nil {
				opts.Logger.Errorf("Can't start watching for metrics: %v", errWatch)
			}
		}()
	} else {
		opts.Logger.Error(err)
	}

	e := echo.New()

	e.Debug = opts.Debug
	e.HideBanner = true
	e.Logger = opts.Logger

	if opts.Validator != nil {
		e.Validator = opts.Validator
	}

	e.HTTPErrorHandler = opts.Error.Capture
	e.Use(opts.Error.Recover())

	return e
}

// StartServer HTTP with custom address.
func StartServer(e *echo.Echo, bind string) error {
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
