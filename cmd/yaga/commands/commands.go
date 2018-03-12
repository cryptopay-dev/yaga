package commands

import (
	"github.com/cryptopay-dev/yaga/config"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli"
	"go.uber.org/atomic"
)

var (
	clr = color.New()
	cnt = atomic.NewInt64(0)
)

var (
	// DefaultLogger for commands:
	DefaultLogger = NewLogger()

	defaultDB   *config.Database
	defaultPath = ""
)

// All returns all commands
func All() cli.Commands {
	clr.Enable()

	log := DefaultLogger

	return []cli.Command{
		newProject(log),                // Creates new project..
		MigrateCreate(defaultPath),     // Creates new migration
		MigrateUp(defaultDB, log),      // Migrations Up to latest
		MigrateDown(defaultDB, log),    // Migrations Down to latest
		MigrateVersion(defaultDB, log), // Get migrations version
		MigrateList(defaultDB, log),    // List applied migrations
		MigratePlan(defaultDB, log),    // Plan to apply migrations
	}
}
