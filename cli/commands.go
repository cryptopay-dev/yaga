package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cryptopay-dev/yaga/migrate"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

func shutdownApplication(opts *Options) {
	if opts.App == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := opts.App.Shutdown(ctx); err != nil {
		opts.Logger.Error(zap.Error(err))
	}
}

func setDatabaseConnector(opts *Options) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		if err = setDatabase(opts, ctx.String("db")); err != nil {
			return err
		}

		if opts.DB == nil {
			return errors.New("database is undefined")
		}

		return nil
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
			// Create context
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			go func() {
				// Stopping server
				if err := opts.App.Shutdown(ctx); err != nil {
					opts.Logger.Fatal("Error stopping server", zap.Error(err))
				}
			}()

			// Running main server
			if err := opts.App.Run(RunOptions{
				DB:           opts.DB,
				Redis:        opts.Redis,
				Logger:       opts.Logger,
				Debug:        opts.Debug,
				BuildTime:    opts.BuildTime,
				BuildVersion: opts.BuildVersion,
			}); err != nil {
				opts.Logger.Fatal("Application failure", zap.Error(err))
			}

			opts.Logger.Info("Application stopped")
			return nil
		},
	})
}

func dbCommands(opts *Options) {
	if opts.enableMigration {
		opts.commands = append(opts.commands, dbCommandSlice(opts)...)
	}
}

func dbCommandSlice(opts *Options) []Command {
	var (
		setNameFlag = cli.StringFlag{
			Name:  "name",
			Value: "",
			Usage: "migration name",
		}
		setStepsFlag = cli.IntFlag{
			Name:  "steps",
			Value: 1,
			Usage: "steps count to down",
		}
		setDBFlag = cli.StringFlag{
			Name:  "db",
			Usage: "set database",
		}
		requiredDBFlag = cli.StringFlag{
			Name:  "db",
			Usage: "set database",
		}
	)

	return cli.Commands{
		{
			Name:  "db:cleanup",
			Usage: "Cleanup database",
			Flags: []cli.Flag{
				requiredDBFlag,
			},
			Before: setDatabaseConnector(opts),
			After: func(context *cli.Context) error {
				shutdownApplication(opts)
				return nil
			},
			Action: func(c *cli.Context) error {
				var (
					names  []string
					dbName = c.String("db")
					db     = opts.DB
				)

				if len(dbName) == 0 {
					return errors.New("you need to set database name `--db <name>`")
				}

				var (
					querySelect   = `SELECT table_name as name FROM information_schema.tables WHERE table_schema = 'public' AND table_name != 'migrations' ORDER BY table_name;`
					queryTruncate = `TRUNCATE %s RESTART IDENTITY;`
				)

				if _, err := db.Query(&names, querySelect); err != nil {
					return err
				}

				if _, err := db.Exec(fmt.Sprintf(
					queryTruncate,
					strings.Join(names, ", "),
				)); err != nil {
					return err
				}

				return nil
			},
		},

		{
			Name:  "migrate:up",
			Usage: "Apply migrations, by default all to newest",
			Flags: []cli.Flag{
				setDBFlag,
			},
			Before: setDatabaseConnector(opts),
			After: func(context *cli.Context) error {
				shutdownApplication(opts)
				return nil
			},
			Action: func(c *cli.Context) error {
				var (
					err      error
					migrator migrate.Migrator
				)

				if migrator, err = migrate.New(migrate.Options{
					DB:     opts.DB,
					Path:   opts.migrationPath,
					Logger: opts.Logger,
				}); err != nil {
					return err
				}

				if err := migrator.Up(0); err != nil {
					opts.Logger.Fatal("Migration failure", zap.Error(err))
				}

				return nil
			},
		},

		{
			Name:  "migrate:down",
			Usage: "Rollback migration by default one",
			Flags: []cli.Flag{
				setDBFlag,
				setStepsFlag,
			},
			Before: setDatabaseConnector(opts),
			After: func(context *cli.Context) error {
				shutdownApplication(opts)
				return nil
			},
			Action: func(c *cli.Context) error {
				var (
					err      error
					steps    = 1
					migrator migrate.Migrator
				)

				if c.Int("steps") > 0 {
					steps = c.Int("steps")
				}

				if migrator, err = migrate.New(migrate.Options{
					DB:     opts.DB,
					Path:   opts.migrationPath,
					Logger: opts.Logger,
				}); err != nil {
					return err
				}

				if err := migrator.Down(steps); err != nil {
					opts.Logger.Fatal("Migration failure", zap.Error(err))
				}

				return nil
			},
		},

		{
			Name:   "migrate:version",
			Usage:  "Current migration version",
			Before: setDatabaseConnector(opts),
			After: func(context *cli.Context) error {
				shutdownApplication(opts)
				return nil
			},
			Action: func(c *cli.Context) error {
				var (
					err      error
					migrator migrate.Migrator
					version  int64
				)

				if migrator, err = migrate.New(migrate.Options{
					DB:     opts.DB,
					Path:   opts.migrationPath,
					Logger: opts.Logger,
				}); err != nil {
					return err
				}

				if version, err = migrator.Version(); err != nil {
					opts.Logger.Fatal("Migration failure", zap.Error(err))
				}

				opts.Logger.Infof("Current version %d", version)

				return nil
			},
		},

		{
			Name:   "migrate:list",
			Usage:  "List current migrations state",
			Before: setDatabaseConnector(opts),
			After: func(context *cli.Context) error {
				shutdownApplication(opts)
				return nil
			},
			Action: func(c *cli.Context) error {
				var (
					err      error
					migrator migrate.Migrator
					items    migrate.Migrations
				)

				if migrator, err = migrate.New(migrate.Options{
					DB:     opts.DB,
					Path:   opts.migrationPath,
					Logger: opts.Logger,
				}); err != nil {
					return err
				}

				if items, err = migrator.List(); err != nil {
					opts.Logger.Fatal("Migration failure", zap.Error(err))
				}

				for _, item := range items {
					opts.Logger.Infof(
						"%s -> %s",
						item.RealName(),
						item.CreatedAt,
					)
				}

				return nil
			},
		},

		{
			Name:   "migrate:plan",
			Usage:  "Current migrations plan",
			Before: setDatabaseConnector(opts),
			After: func(context *cli.Context) error {
				shutdownApplication(opts)
				return nil
			},
			Action: func(c *cli.Context) error {
				var (
					err      error
					migrator migrate.Migrator
					items    migrate.Migrations
				)

				if migrator, err = migrate.New(migrate.Options{
					DB:     opts.DB,
					Path:   opts.migrationPath,
					Logger: opts.Logger,
				}); err != nil {
					return err
				}

				if items, err = migrator.Plan(); err != nil {
					opts.Logger.Fatal("Migration failure", zap.Error(err))
				}

				for _, item := range items {
					opts.Logger.Infof("%s -> not applied", item.RealName())
				}

				return nil
			},
		},

		{
			Name:  "migrate:create",
			Usage: "Create new migration",
			Flags: []cli.Flag{
				setNameFlag,
			},
			Before: setDatabaseConnector(opts),
			After: func(context *cli.Context) error {
				shutdownApplication(opts)
				return nil
			},
			Action: func(c *cli.Context) error {
				var name = c.String("name")
				return migrate.CreateMigration(opts.migrationPath, name)
			},
		},
	}
}
