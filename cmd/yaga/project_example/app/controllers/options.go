package controllers

import (
	"github.com/cryptopay-dev/yaga/web"
)

// Options of controller
type Options struct {
	Engine       *web.Engine
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

// Engine closure to set field in Options
func Engine(e *web.Engine) Option {
	return func(o *Options) {
		o.Engine = e
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
