package cli

import (
	"context"
	"time"

	"github.com/cryptopay-dev/yaga/cmd/yaga/commands"
	"github.com/cryptopay-dev/yaga/config"
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
			ropts, err := NewRunOptions(opts)
			if err != nil {
				return err
			}

			// Running main server
			if err = opts.App.Run(ropts); err != nil {
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
