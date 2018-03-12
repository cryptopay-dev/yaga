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

	defaultDB   *config.Database // nil
	defaultPath = ""             // empty
)

// All returns all commands
func All() cli.Commands {
	clr.Enable()

	return []cli.Command{
		newProject(DefaultLogger),                // Creates new project..
		MigrateCreate(defaultPath),               // Creates new migration
		MigrateUp(defaultDB, DefaultLogger),      // Migrations Up to latest
		MigrateDown(defaultDB, DefaultLogger),    // Migrations Down to latest
		MigrateVersion(defaultDB, DefaultLogger), // Get migrations version
		MigrateList(defaultDB, DefaultLogger),    // List applied migrations
		MigratePlan(defaultDB, DefaultLogger),    // Plan to apply migrations
		MigrateCleanup(defaultDB, DefaultLogger), // Cleanup database...
	}
}
