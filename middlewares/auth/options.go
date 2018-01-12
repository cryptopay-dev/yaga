package auth

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
)

type Options struct {
	DB     *pg.DB
	Logger logger.Logger
}

type Option func(*Options)

func newOptions(opts ...Option) (opt Options) {
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func DB(db *pg.DB) Option {
	return func(o *Options) {
		o.DB = db
	}
}

func Logger(log logger.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}
