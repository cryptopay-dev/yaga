package commands

import (
	"github.com/urfave/cli"
)

// MigrateUp migrations
func MigrateUp() cli.Command {
	return cli.Command{
		Name:        "migrate:up",
		ShortName:   "m:u",
		Usage:       "up --steps=<count> --dsn=<DSN> --db=<db-name> --path=<to-migrations>",
		Description: "Migration up to latest migration (by default)",
		Category:    "Migrate commands",
		Flags:       migrateFlags(),
		Action:      migrateAction(migrateUp),
	}
}
