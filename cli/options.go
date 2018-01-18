package cli

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
	"github.com/urfave/cli"
	"gopkg.in/go-playground/validator.v9"
)

type Options struct {
	App          Instance      `validate:"required"`
	DB           *pg.DB        `validate:"required"`
	Level        string        `validate:"required"`
	Logger       logger.Logger `validate:"required"`
	Users        []cli.Author
	Usage        string
	Name         string
	BuildTime    string
	BuildVersion string
}

type Option func(*Options)

func newOptions(opts ...Option) (opt Options, err error) {
	for _, o := range opts {
		o(&opt)
	}
	err = validator.New().Struct(opt)
	return
}

func App(app Instance) Option {
	return func(o *Options) {
		o.App = app
	}
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

func BuildTime(buildTime string) Option {
	return func(o *Options) {
		o.BuildTime = buildTime
	}
}

func BuildVersion(buildVersion string) Option {
	return func(o *Options) {
		o.BuildVersion = buildVersion
	}
}

func Level(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func Usage(descr string) Option {
	return func(o *Options) {
		o.Usage = descr
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Users(users []cli.Author) Option {
	return func(o *Options) {
		o.Users = users
	}
}
