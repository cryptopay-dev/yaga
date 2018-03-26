package commands

import (
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

	defaultPath = "" // empty
)

// All returns all commands
func All() cli.Commands {
	clr.Enable()

	return []cli.Command{
		newProject(DefaultLogger),     // Creates new project..
		MigrateCreate(defaultPath),    // Creates new migration
		MigrateUp(DefaultLogger),      // Migrations Up to latest
		MigrateDown(DefaultLogger),    // Migrations Down to latest
		MigrateVersion(DefaultLogger), // Get migrations version
		MigrateList(DefaultLogger),    // List applied migrations
		MigratePlan(DefaultLogger),    // Plan to apply migrations
		MigrateCleanup(DefaultLogger), // Cleanup database...
	}
}
