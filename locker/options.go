package locker

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-redis/redis"
)

type Options struct {
	Redis  *redis.Client
	Logger logger.Logger
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	return options
}

func Redis(r *redis.Client) Option {
	return func(o *Options) {
		o.Redis = r
	}
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}
