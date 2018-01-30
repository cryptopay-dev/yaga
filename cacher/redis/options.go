package redis

// Options for creating redis-cacher
type Options struct {
	Address  string
	Password string
	DB       int
}

// newOptions converts closure to Options-field
func newOptions(opts ...Option) Options {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Option closure
type Option func(*Options)

// Address closure to set in Options
func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// Password closure to set in Options
func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

// DB closure to set in Options
func DB(db int) Option {
	return func(o *Options) {
		o.DB = db
	}
}
