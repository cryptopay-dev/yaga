package auth

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
)

// Options for auth
type Options struct {
	DB     *pg.DB
	Logger logger.Logger
}

// Option closure
type Option func(*Options)

// newOptions sets values to Options
func newOptions(opts ...Option) (opt Options) {
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// DB closure to set value in Options
func DB(db *pg.DB) Option {
	return func(o *Options) {
		o.DB = db
	}
}

// Logger closure to set value in Options
func Logger(log logger.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}
