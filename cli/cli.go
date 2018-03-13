package cli

import (
	"errors"
	"os"
	"reflect"
	"sort"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/cryptopay-dev/yaga/validate"
	wrap "github.com/pkg/errors"
	"github.com/urfave/cli"
	"gopkg.in/go-playground/validator.v9"
)

var (
	// ErrAppNotPointer when app-instance not pointer to struct
	ErrAppNotPointer = errors.New("app must be a pointer to a struct")
	// ErrConfigNotPointer when app-instance not pointer to struct
	ErrConfigNotPointer = errors.New("config must be a pointer to a struct")
)

// Run creates instance of cli.App with Options.
// Validate options with https://github.com/go-playground/validator
// Required:
// - App instance
// - Logger
func Run(opts ...Option) error {
	var options = newOptions(opts...)

	cliApp := cli.NewApp()
	cliApp.Name = options.Name
	cliApp.Usage = options.Usage
	cliApp.Version = options.BuildVersion
	cliApp.Authors = options.Users

	cliApp.Before = before(options)

	if options.action != nil {
		cliApp.Action = options.action
	}
	if options.after != nil {
		cliApp.After = options.after
	}
	if len(options.flags) > 0 {
		cliApp.Flags = append(cliApp.Flags, options.flags...)
	}

	appCommands(options)
	dbCommands(options)

	if len(options.commands) > 0 {
		cliApp.Commands = append(cliApp.Commands, options.commands...)
	}

	sort.Sort(cli.CommandsByName(cliApp.Commands))

	return cliApp.Run(os.Args)
}

func before(options *Options) func(ctx *Context) error {
	return func(ctx *Context) (err error) {
		if options.before != nil {
			if err = options.before(ctx); err != nil {
				return err
			}
		}

		if options.Logger == nil {
			if options.Debug == false { // Debug = false
				options.Logger = zap.New(zap.Production)
			} else if options.Quiet { // Debug = true && Quiet = true
				options.Logger = nop.New()
			} else { // Debug = true && Quiet = false
				options.Logger = zap.New(zap.Development)
			}
		}

		// If we have config-source/interface - loading config:
		if options.ConfigSource != nil &&
			options.ConfigInterface != nil {
			if reflect.TypeOf(options.ConfigInterface).Kind() != reflect.Ptr {
				return ErrConfigNotPointer
			}

			if err = config.Load(
				options.ConfigSource,
				options.ConfigInterface,
			); err != nil {
				return err
			}
		}

		if options.App != nil && reflect.TypeOf(options.App).Kind() != reflect.Ptr {
			return ErrAppNotPointer
		}

		if err = setDatabase(options, ""); err != nil {
			return wrap.Wrap(err, "can't connect to database")
		}

		if options.ConfigInterface != nil {
			if redisConf, ok := hasRedis(options.ConfigInterface); ok {
				if options.Redis, err = redisConf.Connect(); err != nil {
					return wrap.Wrap(err, "can't connect to redis")
				}
			}
		}

		// Validate options:
		if err = validator.New().Struct(options); err != nil {
			if ok, errVal := validate.CheckErrors(validate.Options{
				Struct: options,
				Errors: err,
			}); ok {
				return wrap.Wrap(errVal, "options not valid!")
			}
		}

		return err
	}
}

func setDatabase(opts *Options, dbname string) (err error) {
	if (opts.DB != nil && len(dbname) == 0) || opts.ConfigInterface == nil {
		return nil
	}

	dbConf, ok := hasDB(opts.ConfigInterface)
	if !ok {
		// TODO or return an error?
		return nil
	}

	if len(dbname) != 0 {
		dbConf.Database = dbname
	}

	opts.DB, err = dbConf.Connect()

	return err
}

func hasRedis(conf interface{}) (*config.Redis, bool) {
	v := reflect.ValueOf(conf).Elem()

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		if val, ok := v.Field(i).Interface().(config.Redis); ok {
			return &val, true
		}
	}

	return nil, false
}

func hasDB(conf interface{}) (*config.Database, bool) {
	v := reflect.ValueOf(conf).Elem()

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}

		if val, ok := v.Field(i).Interface().(config.Database); ok {
			return &val, true
		}
	}

	return nil, false
}
