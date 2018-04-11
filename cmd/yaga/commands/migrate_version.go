package commands

import (
	"github.com/urfave/cli"
)

// MigrateVersion migrations
func MigrateVersion() cli.Command {
	return cli.Command{
		Name:        "migrate:version",
		ShortName:   "m:v",
		Usage:       "version --db=<db-name> --dsn=<DSN>",
		Description: "Migration version",
		Category:    "Migrate commands",
		Flags:       []cli.Flag{dbFlag, dsnFlag},
		Action:      migrateAction(migrateVersion),
	}
}
