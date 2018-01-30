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

type Options struct {
	Logger logger.Logger
	Error  errors.Logic
	Debug  bool
}

type Context = echo.Context

func New(opts Options) *echo.Echo {
	// Enabling raven:
	if err := raven.SetDSN(os.Getenv("SENTRY_DSN")); err != nil {
		opts.Logger.Error(err)
	}

	//Enable metrics:
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

	e.HTTPErrorHandler = opts.Error.Capture
	e.Use(opts.Error.Recover())

	return e
}

func StartServer(e *echo.Echo, bind string) error {
	// Start server
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
