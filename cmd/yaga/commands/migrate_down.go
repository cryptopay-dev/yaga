package commands

import (
	"github.com/urfave/cli"
)

// MigrateDown migrations
func MigrateDown() cli.Command {
	return cli.Command{
		Name:        "migrate:down",
		ShortName:   "m:d",
		Usage:       "down --steps=<count> --dsn=<DSN> --db=<db-name> --path=<to-migrations>",
		Description: "Migration down last migration (by default)",
		Category:    "Migrate commands",
		Flags:       migrateFlags(),
		Action:      migrateAction(migrateDown),
	}
}
