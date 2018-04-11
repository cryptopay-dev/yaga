package cli

import (
	"context"
	"reflect"
	"time"

	"github.com/cryptopay-dev/yaga/cmd/yaga/commands"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/urfave/cli"
)

func shutdownApplication(opts *Options) {
	if opts.App == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := opts.App.Shutdown(ctx); err != nil {
		log.Error(err)
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
			log.Info("Application stopped")
			return nil
		},
		Action: func(c *cli.Context) error {
			if opts.App != nil && reflect.TypeOf(opts.App).Kind() != reflect.Ptr {
				return ErrAppNotPointer
			}

			// Running main server
			return opts.App.Run(RunOptions{
				Debug:        opts.Debug,
				BuildTime:    opts.BuildTime,
				BuildVersion: opts.BuildVersion,
			})
		},
	})
}

func dbCommands(opts *Options) {
	opts.commands = append(opts.commands, dbCommandSlice(opts)...)
}

func dbCommandSlice(opts *Options) []Command {
	return cli.Commands{
		// Migrate cleanup
		commands.MigrateCleanup(),

		// Migrate up
		commands.MigrateUp(),

		// Migrate down
		commands.MigrateDown(),

		// Migrate version:
		commands.MigrateVersion(),

		// List applied migrations:
		commands.MigrateList(),

		// List plan to migrate:
		commands.MigratePlan(),

		// Create migrations:
		commands.MigrateCreate(opts.migrationPath),
	}
}
