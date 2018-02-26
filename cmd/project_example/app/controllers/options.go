package controllers

import (
	"github.com/cryptopay-dev/yaga/cmd/project_example/app/library/config"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/cryptopay-dev/yaga/web"
)

// Options of controller
type Options struct {
	Conf         *config.Config
	Engine       *web.Engine
	Logger       logger.Logger
	BuildTime    string
	BuildVersion string
}

// Option closure
type Option func(*Options)

// newOptions transform closure to Options
func newOptions(opts ...Option) (opt Options) {
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Config closure to set field in Options
func Config(c *config.Config) Option {
	return func(o *Options) {
		o.Conf = c
	}
}

// Engine closure to set field in Options
func Engine(e *web.Engine) Option {
	return func(o *Options) {
		o.Engine = e
	}
}

// Logger closure to set field in Options
func Logger(log logger.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}

// BuildTime closure to set field in Options
func BuildTime(buildTime string) Option {
	return func(o *Options) {
		o.BuildTime = buildTime
	}
}

// BuildVersion closure to set field in Options
func BuildVersion(buildVersion string) Option {
	return func(o *Options) {
		o.BuildVersion = buildVersion
	}
}
