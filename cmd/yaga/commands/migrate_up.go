package commands

import (
	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/urfave/cli"
)

// MigrateUp migrations
func MigrateUp(db *config.Database, log logger.Logger) cli.Command {
	return cli.Command{
		Name:        "migrate:up",
		ShortName:   "m:u",
		Usage:       "up --steps=<count> --dsn=<DSN>",
		Description: "Migration up to latest migration (by default)",
		Category:    "Migrate commands",
		Flags:       migrateFlags(),
		Action:      migrateAction(migrateUp, db, log),
	}
}
