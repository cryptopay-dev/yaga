package commands

import (
	"fmt"
	"strings"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/helpers/postgres"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
	"github.com/labstack/gommon/log"
	"github.com/urfave/cli"
)

// MigrateCleanup migrations
func MigrateCleanup(dlog logger.Logger) cli.Command {
	action := func(ctx *cli.Context) (err error) {
		if err = FetchDB(ctx, "database"); err != nil {
			log.Fatalf("can't find config file or dsn: %v", err)
		}

		if database := ctx.String("db"); len(database) != 0 {
			config.Set("database.database", database)
		}

		database, err := postgres.Connect("database")
		if err != nil {
			log.Fatalf("postgres connection error: %v", err)
		}

		database.RunInTransaction(func(tx *pg.Tx) error {
			var (
				names         []string
				querySelect   = `SELECT table_name as name FROM information_schema.tables WHERE table_schema = 'public' AND table_name != 'migrations' ORDER BY table_name;`
				queryTruncate = `TRUNCATE %s RESTART IDENTITY;`
			)

			if _, err := tx.Query(&names, querySelect); err != nil {
				return err
			}

			if _, err := tx.Exec(fmt.Sprintf(
				queryTruncate,
				strings.Join(names, ", "),
			)); err != nil {
				return err
			}

			return nil
		})

		return nil
	}

	return cli.Command{
		Name:        "migrate:cleanup",
		ShortName:   "m:cl",
		Usage:       "cleanup --db=<db-name> --dsn=<DSN>",
		Description: "Migration cleanup",
		Category:    "Migrate commands",
		Flags:       []cli.Flag{dbFlag, dsnFlag},
		Action:      action,
	}
}
