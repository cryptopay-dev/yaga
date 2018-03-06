package controllers

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/cmd/project_example/app/library/config"
	"github.com/cryptopay-dev/yaga/cmd/project_example/misc"
	"github.com/cryptopay-dev/yaga/doc"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/cryptopay-dev/yaga/pprof"
	"github.com/cryptopay-dev/yaga/web"
)

const swaggerFile = "assets/docs/swagger.yaml"

// Controller - rest-api controller
type Controller struct {
	Conf         *config.Config
	Engine       *web.Engine
	Logger       logger.Logger
	BuildTime    string
	BuildVersion string
}

// JSON alias
type JSON = map[string]interface{}

// New - Create new REST-API
func New(opts ...Option) (*Controller, error) {
	var (
		options = newOptions(opts...)
		ctrl    = Controller{
			Conf:         options.Conf,
			Engine:       options.Engine,
			Logger:       options.Logger,
			BuildTime:    options.BuildTime,
			BuildVersion: options.BuildVersion,
		}
	)

	// Debug:
	if err := pprof.Wrap(options.Logger, ctrl.Engine); err != nil {
		return nil, err
	}

	// Version:
	ctrl.Engine.GET("/version", ctrl.Version)

	// Doc
	doc.AddDocumentation(ctrl.Engine, "", misc.Name, swaggerFile)

	return &ctrl, nil
}

// Version of application
func (c *Controller) Version(ctx web.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"time":    c.BuildTime,
		"version": c.BuildVersion,
	})
}
