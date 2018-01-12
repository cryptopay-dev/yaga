package redis

type Options struct {
	Address  string
	Password string
	DB       int
}

func newOptions(opts ...Option) Options {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	return options
}

type Option func(*Options)

func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

func Password(password string) Option {
	return func(o *Options) {
		o.Password = password
	}
}

func DB(db int) Option {
	return func(o *Options) {
		o.DB = db
	}
}
