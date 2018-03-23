package commands

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/urfave/cli"
)

// MigrateDown migrations
func MigrateDown(log logger.Logger) cli.Command {
	return cli.Command{
		Name:        "migrate:down",
		ShortName:   "m:d",
		Usage:       "down --steps=<count> --dsn=<DSN> --db=<db-name> --path=<to-migrations>",
		Description: "Migration down last migration (by default)",
		Category:    "Migrate commands",
		Flags:       migrateFlags(),
		Action:      migrateAction(migrateDown, log),
	}
}
