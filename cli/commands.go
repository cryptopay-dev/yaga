package cli

import (
	"context"
	"reflect"
	"time"

	"github.com/cryptopay-dev/yaga/cmd/yaga/commands"
	"github.com/urfave/cli"
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

			if opts.App != nil && reflect.TypeOf(opts.App).Kind() != reflect.Ptr {
				return ErrAppNotPointer
			}

			// Running main server
			if err = opts.App.Run(RunOptions{
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
	return cli.Commands{
		// Migrate cleanup
		commands.MigrateCleanup(opts.Logger),

		// Migrate up
		commands.MigrateUp(opts.Logger),

		// Migrate down
		commands.MigrateDown(opts.Logger),

		// Migrate version:
		commands.MigrateVersion(opts.Logger),

		// List applied migrations:
		commands.MigrateList(opts.Logger),

		// List plan to migrate:
		commands.MigratePlan(opts.Logger),

		// Create migrations:
		commands.MigrateCreate(opts.migrationPath),
	}
}
