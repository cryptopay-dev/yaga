package cli

import (
	"context"

	"github.com/urfave/cli"
)

// RunOptions for pass db, redis, etc to application:
type RunOptions struct {
	Debug        bool
	BuildTime    string
	BuildVersion string
}

// Instance abstraction layer above Application
type Instance interface {
	Run(RunOptions) error
	Shutdown(ctx context.Context) error
}

type (
	// Flag is a common interface related to parsing flags in cli.
	// For more advanced flag parsing techniques, it is recommended that
	// this interface be implemented.
	Flag = cli.Flag
	// IntFlag is a flag with type int
	IntFlag = cli.IntFlag
	// StringFlag is a flag with type string
	StringFlag = cli.StringFlag
	// Author represents someone who has contributed to a cli project.
	Author = cli.Author
	// Context is a type that is passed through to
	// each Handler action in a cli application. Context
	// can be used to retrieve context-specific Args and
	// parsed command-line options.
	Context = cli.Context
	// Command is a subcommand for a cli.App.
	Command = cli.Command
	// Commandor closure for applying options to command
	Commandor func(*Options) Command
)
