package commands

import (
	"github.com/urfave/cli"
)

var defaultPath = "" // empty

// All returns all commands
func All() cli.Commands {
	return []cli.Command{
		newProject(),               // Creates new project..
		MigrateCreate(defaultPath), // Creates new migration
		MigrateUp(),                // Migrations Up to latest
		MigrateDown(),              // Migrations Down to latest
		MigrateVersion(),           // Get migrations version
		MigrateList(),              // List applied migrations
		MigratePlan(),              // Plan to apply migrations
		MigrateCleanup(),           // Cleanup database...
	}
}
