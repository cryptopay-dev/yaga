package cli

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
	"github.com/urfave/cli"
)

// Options for creating cli.App instance
type Options struct {
	App             Instance
	Logger          logger.Logger `validate:"required"`
	ConfigSource    interface{}   `validate:"required"`
	ConfigInterface interface{}   `validate:"required"`
	DB              *pg.DB
	Redis           *redis.Client
	Users           []cli.Author
	Debug           bool
	Quiet           bool
	Usage           string
	Name            string
	BuildTime       string
	BuildVersion    string

	action        func(*Context) error
	before        func(*Context) error
	after         func(*Context) error
	commands      []Command
	flags         []Flag
	migrationPath string
}

// Option closure
type Option func(*Options)

// newOptions converts slice of closures to Options-field
func newOptions(opts ...Option) (opt *Options) {
	opt = &Options{
		// TODO need?
		migrationPath: "./migrations",
	}

	for _, o := range opts {
		o(opt)
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

// Config closure to set config source and interface in Options
func Config(src, conf interface{}) Option {
	return func(o *Options) {
		o.ConfigSource = src
		o.ConfigInterface = conf
	}
}

// MigrationsPath closure to set param in Options
func MigrationsPath(path string) Option {
	return func(o *Options) {
		o.migrationPath = path
	}
}

// Debug closure to set debug and quiet state of logger in Options
func Debug(args ...bool) Option {
	return func(o *Options) {
		if len(args) > 0 {
			o.Debug = args[0]
		}

		if len(args) > 1 {
			o.Quiet = args[1]
		}
	}
}

// Commands closure to set additional commands for CLI
func Commands(cmds ...Commandor) Option {
	return func(o *Options) {
		o.commands = make([]Command, 0, len(cmds))

		for _, cmd := range cmds {
			o.commands = append(o.commands, cmd(o))
		}
	}
}

// Flags closure to set additional commands for CLI
func Flags(flags ...Flager) Option {
	return func(o *Options) {
		o.flags = make([]Flag, 0, len(flags))

		for _, flag := range flags {
			o.flags = append(o.flags, flag(o))
		}
	}
}

// Trigger closure to set triggers for CLI
func Trigger(action, before, after Handler) Option {
	return func(o *Options) {
		if action != nil {
			o.action = action(o)
		}
		if before != nil {
			o.before = before(o)
		}
		if after != nil {
			o.after = after(o)
		}
	}
}
