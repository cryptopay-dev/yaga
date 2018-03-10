package commands

import (
	"os"
	"path"
	"strings"

	"github.com/cryptopay-dev/yaga/migrate"
	"github.com/urfave/cli"
)

const defaultMigrationsPath = "migrations"

// MigrateCreate new migration files
func MigrateCreate(defaultPath string) cli.Command {
	if len(defaultPath) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			dir = "./"
		}

		defaultPath = path.Join(dir, defaultMigrationsPath)
	}

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "migration path",
			Value: defaultPath,
		},
	}

	action := func(ctx *cli.Context) error {
		mpath := ctx.String("path")

		if _, err := os.Stat(mpath); err != nil {
			errorsf("migration path not found: %v", err)
		}

		name := strings.Join(ctx.Args(), "_")
		name = strings.ToLower(name)

		if len(name) == 0 {
			errors("migration name not set")
		}

		if err := migrate.CreateMigration(mpath, name); err != nil {
			errorsf("migration not created: %v", err)
		}

		infof("migration created: %s", name)

		return nil
	}

	return cli.Command{
		Name:        "migrate:create",
		ShortName:   "n",
		Usage:       "new <migration-name>",
		Description: "Create new migration",
		Category:    "migrate",
		Flags:       flags,
		Action:      action,
	}
}
