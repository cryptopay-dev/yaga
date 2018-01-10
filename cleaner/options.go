package cleaner

import (
	"github.com/cryptopay-dev/yaga/cacher"
	"github.com/cryptopay-dev/yaga/logger"
)

type Options struct {
	Platforms Platforms
	Logger    logger.Logger
	Cacher    cacher.Cacher
}

type Option func(opts *Options)

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Platform(platforms Platforms) Option {
	return func(o *Options) {
		o.Platforms = platforms
	}
}

func Logger(log logger.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}

func Cacher(cache cacher.Cacher) Option {
	return func(o *Options) {
		o.Cacher = cache
	}
}
