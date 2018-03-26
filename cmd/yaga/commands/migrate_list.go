package commands

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/urfave/cli"
)

// MigrateList migrations
func MigrateList(log logger.Logger) cli.Command {
	return cli.Command{
		Name:        "migrate:list",
		ShortName:   "m:l",
		Usage:       "list --db=<db-name> --dsn=<DSN>",
		Description: "Migration list applied migrations",
		Category:    "Migrate commands",
		Flags:       []cli.Flag{dbFlag, dsnFlag},
		Action:      migrateAction(migrateList, log),
	}
}

// MigratePlan migrations
func MigratePlan(log logger.Logger) cli.Command {
	return cli.Command{
		Name:        "migrate:plan",
		ShortName:   "m:p",
		Usage:       "plan --db=<db-name> --dsn=<DSN> --db=<db-name> --path=<to-migrations>",
		Description: "Migration plan migrations",
		Category:    "Migrate commands",
		Flags:       migrateFlags(),
		Action:      migrateAction(migratePlan, log),
	}
}
