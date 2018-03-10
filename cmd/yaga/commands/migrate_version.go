package commands

import (
	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/urfave/cli"
)

// MigrateVersion migrations
func MigrateVersion(db *config.Database, log logger.Logger) cli.Command {
	return cli.Command{
		Name:        "migrate:version",
		ShortName:   "m:v",
		Usage:       "version --db=<db-name> --dsn=<DSN>",
		Description: "Migration version",
		Category:    "Migrate commands",
		Flags:       []cli.Flag{dbFlag, dsnFlag},
		Action:      migrateAction(migrateVersion, db, log),
	}
}
