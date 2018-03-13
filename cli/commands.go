package cli

import (
	"context"
	"reflect"
	"time"

	"github.com/cryptopay-dev/yaga/cmd/yaga/commands"
	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/validate"
	"github.com/urfave/cli"
	"gopkg.in/go-playground/validator.v9"
)

func shutdownApplication(opts *Options) {
	if opts.App == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := opts.App.Shutdown(ctx); err != nil {
		opts.Logger.Error(err)
	}
}

func appCommands(opts *Options) {
	if opts.App == nil {
		return
	}

	opts.commands = append(opts.commands, Command{
		Name:    "start",
		Aliases: []string{"s"},
		Usage:   "start main server",
		After: func(context *cli.Context) error {
			shutdownApplication(opts)
			return nil
		},
		Action: func(c *cli.Context) error {
			var err error

			// If we have config-source/interface - loading config:
			if opts.ConfigSource != nil &&
				opts.ConfigInterface != nil {
				if reflect.TypeOf(opts.ConfigInterface).Kind() != reflect.Ptr {
					return ErrConfigNotPointer
				}

				if err = config.Load(
					opts.ConfigSource,
					opts.ConfigInterface,
				); err != nil {
					return err
				}
			}

			if opts.App != nil && reflect.TypeOf(opts.App).Kind() != reflect.Ptr {
				return ErrAppNotPointer
			}

			if err = setDatabase(opts, ""); err != nil {
				return err
			}

			if opts.ConfigInterface != nil {
				if redisConf, ok := hasRedis(opts.ConfigInterface); ok {
					if opts.Redis, err = redisConf.Connect(); err != nil {
						return err
					}
				}
			}

			// Validate options:
			if err := validator.New().Struct(opts); err != nil {
				if ok, errv := validate.CheckErrors(validate.Options{
					Struct: opts,
					Errors: err,
				}); ok {
					panic(errv)
				}
			}

			// Running main server
			if err := opts.App.Run(RunOptions{
				DB:           opts.DB,
				Redis:        opts.Redis,
				Logger:       opts.Logger,
				Debug:        opts.Debug,
				BuildTime:    opts.BuildTime,
				BuildVersion: opts.BuildVersion,
			}); err != nil {
				opts.Logger.Fatal("Application failure", err)
			}

			opts.Logger.Info("Application stopped")
			return nil
		},
	})
}

func dbCommands(opts *Options) {
	opts.commands = append(opts.commands, dbCommandSlice(opts)...)
}

func dbCommandSlice(opts *Options) []Command {
	var db *config.Database

	if opts.DB != nil {
		conf := opts.DB.Options()
		db = &config.Database{
			Address:  conf.Addr,
			Database: conf.Database,
			User:     conf.User,
			Password: conf.Password,
		}
	}

	return cli.Commands{
		// Migrate cleanup
		commands.MigrateCleanup(db, opts.Logger),

		// Migrate up
		commands.MigrateUp(db, opts.Logger),

		// Migrate down
		commands.MigrateDown(db, opts.Logger),

		// Migrate version:
		commands.MigrateVersion(db, opts.Logger),

		// List applied migrations:
		commands.MigrateList(db, opts.Logger),

		// List plan to migrate:
		commands.MigratePlan(db, opts.Logger),

		// Create migrations:
		commands.MigrateCreate(opts.migrationPath),
	}
}
