package auth

import (
	"github.com/go-pg/pg"
)

// Options for auth
type Options struct {
	DB *pg.DB
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
