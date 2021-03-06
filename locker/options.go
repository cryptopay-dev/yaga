package locker

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-redis/redis"
)

// Options for creating Locker instance
type Options struct {
	Redis  *redis.Client
	Logger logger.Logger
}

// Option closure
type Option func(*Options)

// newOptions converts slice of closures to Options-field
func newOptions(opts ...Option) Options {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Redis closure to set field in Options
func Redis(r *redis.Client) Option {
	return func(o *Options) {
		o.Redis = r
	}
}

// Logger closure to set field in Options
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}
