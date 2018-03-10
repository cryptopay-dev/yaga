package commands

import (
	"os"
	"path"
	"strings"

	"github.com/cryptopay-dev/yaga/migrate"
	"github.com/labstack/gommon/log"
	"github.com/urfave/cli"
)

const defaultMigrationsPath = "migrations"

func defaultMigratePath() string {
	dir, err := os.Getwd()
	if err != nil {
		dir = "./"
	}

	return path.Join(dir, defaultMigrationsPath)
}

// MigrateCreate new migration files
func MigrateCreate(defaultPath string) cli.Command {
	if len(defaultPath) == 0 {
		defaultPath = defaultMigratePath()
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
			log.Fatalf("migration path not found: %v", err)
		}

		name := strings.Join(ctx.Args(), "_")
		name = strings.ToLower(name)

		if len(name) == 0 {
			log.Fatal("migration name not set")
		}

		if err := migrate.CreateMigration(mpath, name); err != nil {
			log.Fatalf("migration not created: %v", err)
		}

		log.Infof("migration created: %s", name)

		return nil
	}

	return cli.Command{
		Name:        "migrate:create",
		ShortName:   "m:c",
		Usage:       "new <migration-name>",
		Description: "Create new migration",
		Category:    "Migrate commands",
		Flags:       flags,
		Action:      action,
	}
}
