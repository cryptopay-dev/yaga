package commands

import (
	"github.com/urfave/cli"
)

// MigrateList migrations
func MigrateList() cli.Command {
	return cli.Command{
		Name:        "migrate:list",
		ShortName:   "m:l",
		Usage:       "list --db=<db-name> --dsn=<DSN>",
		Description: "Migration list applied migrations",
		Category:    "Migrate commands",
		Flags:       []cli.Flag{dbFlag, dsnFlag},
		Action:      migrateAction(migrateList),
	}
}

// MigratePlan migrations
func MigratePlan() cli.Command {
	return cli.Command{
		Name:        "migrate:plan",
		ShortName:   "m:p",
		Usage:       "plan --db=<db-name> --dsn=<DSN> --db=<db-name> --path=<to-migrations>",
		Description: "Migration plan migrations",
		Category:    "Migrate commands",
		Flags:       migrateFlags(),
		Action:      migrateAction(migratePlan),
	}
}
