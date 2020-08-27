package cli

import (
	"reflect"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/validate"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

func NewRunOptions(opts *Options) (RunOptions, error) {
	var (
		err   error
		ropts RunOptions
	)

	// If we have config-source/interface - loading config:
	if opts.ConfigSource != nil &&
		opts.ConfigInterface != nil {
		if reflect.TypeOf(opts.ConfigInterface).Kind() != reflect.Ptr {
			return ropts, ErrConfigNotPointer
		}

		if err = config.Load(
			opts.ConfigSource,
			opts.ConfigInterface,
		); err != nil {
			return ropts, errors.Wrapf(err, "can't load config")
		}
	}

	if opts.App != nil && reflect.TypeOf(opts.App).Kind() != reflect.Ptr {
		return ropts, ErrAppNotPointer
	}

	if err = setDatabase(opts, ""); err != nil {
		return ropts, errors.Wrapf(err, "can't set database")
	}

	if opts.ConfigInterface != nil {
		if redisConf, ok := hasRedis(opts.ConfigInterface); ok {
			if opts.Redis, err = redisConf.Connect(); err != nil {
				return ropts, errors.Wrap(err, "can't connect to redis")
			}
		}
	}

	// Validate options:
	if err = validator.New().Struct(opts); err != nil {
		if ok, errv := validate.CheckErrors(validate.Options{
			Struct: opts,
			Errors: err,
		}); ok {
			return ropts, errors.Wrapf(errv, "validate error")
		}
	}

	return RunOptions{
		DB:           opts.DB,
		Redis:        opts.Redis,
		Logger:       opts.Logger,
		Debug:        opts.Debug,
		BuildTime:    opts.BuildTime,
		BuildVersion: opts.BuildVersion,
	}, nil
}
