package cli

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
	"github.com/urfave/cli"
)

// Options for creating cli.App instance
type Options struct {
	App             Instance      `validate:"required"`
	Logger          logger.Logger `validate:"required"`
	ConfigSource    interface{}   `validate:"required"`
	ConfigInterface interface{}   `validate:"required"`
	DB              *pg.DB
	Redis           *redis.Client
	Users           []cli.Author
	Usage           string
	Name            string
	BuildTime       string
	BuildVersion    string
}

// Option closure
type Option func(*Options)

// newOptions converts slice of closures to Options-field
func newOptions(opts ...Option) (opt Options) {
	for _, o := range opts {
		o(&opt)
	}
	return
}

// App closure to set field in Options
func App(app Instance) Option {
	return func(o *Options) {
		o.App = app
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

// Usage closure to set field in Options
func Usage(descr string) Option {
	return func(o *Options) {
		o.Usage = descr
	}
}

// Name closure to set field in Options
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Users closure to set field in Options
func Users(users []cli.Author) Option {
	return func(o *Options) {
		o.Users = users
	}
}

// ConfigSource interface for loading config
func ConfigSource(v interface{}) Option {
	return func(o *Options) {
		o.ConfigSource = v
	}
}

// ConfigInterface for loading config
func ConfigInterface(v interface{}) Option {
	return func(o *Options) {
		o.ConfigInterface = v
	}
}
