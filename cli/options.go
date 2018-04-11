package cli

import (
	"github.com/urfave/cli"
)

// Options for creating cli.App instance
type Options struct {
	App          Instance
	Users        []cli.Author
	Debug        bool
	Quiet        bool
	Usage        string
	Name         string
	BuildTime    string
	BuildVersion string

	commands      []Command
	flags         []Flag
	migrationPath string
}

// Option closure
type Option func(*Options)

// newOptions converts slice of closures to Options-field
func newOptions(opts ...Option) (opt *Options) {
	opt = &Options{
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
func Commands(commands ...Commandor) Option {
	return func(o *Options) {
		o.commands = make([]Command, 0, len(commands))

		for _, cmd := range commands {
			o.commands = append(o.commands, cmd(o))
		}
	}
}

// Flags closure to set additional commands for CLI
func Flags(flags ...Flag) Option {
	return func(o *Options) {
		o.flags = flags
	}
}
