package cli

import (
	"errors"
	"os"
	"reflect"
	"sort"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/logger/zap"
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
	var (
		err     error
		options = newOptions(opts...)
	)

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

	if reflect.TypeOf(options.App).Kind() != reflect.Ptr {
		return ErrAppNotPointer
	}

	// Validate options:
	if err = validator.New().Struct(options); err != nil {
		return err
	}

	cliApp := cli.NewApp()
	cliApp.Name = options.Name
	cliApp.Usage = options.Usage
	cliApp.Version = options.BuildVersion
	cliApp.Authors = options.Users

	if dbConf, ok := hasDB(options.ConfigInterface); ok {
		if options.DB, err = dbConf.Connect(); err != nil {
			return err
		}
	}

	if redisConf, ok := hasRedis(options.ConfigInterface); ok {
		if options.Redis, err = redisConf.Connect(); err != nil {
			return err
		}
	}

	addCommands(cliApp, options)
	sort.Sort(cli.CommandsByName(cliApp.Commands))

	return cliApp.Run(os.Args)
}

func hasRedis(conf interface{}) (*config.Redis, bool) {
	v := reflect.ValueOf(conf).Elem()

	for i := 0; i < v.NumField(); i++ {
		if val, ok := v.Field(i).Interface().(config.Redis); ok {
			return &val, true
		}
	}

	return nil, false
}

func hasDB(conf interface{}) (*config.Database, bool) {
	v := reflect.ValueOf(conf).Elem()

	for i := 0; i < v.NumField(); i++ {
		if val, ok := v.Field(i).Interface().(config.Database); ok {
			return &val, true
		}
	}

	return nil, false
}
