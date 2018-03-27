package cli

import (
	"errors"
	"os"
	"sort"

	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/urfave/cli"
)

var (
	// ErrAppNotPointer when app-instance not pointer to struct
	ErrAppNotPointer = errors.New("app must be a pointer to a struct")
)

// Run creates instance of cli.App with Options.
// Validate options with https://github.com/go-playground/validator
// Required:
// - App instance
// - Logger
func Run(opts ...Option) error {
	var options = newOptions(opts...)

	cliApp := cli.NewApp()
	cliApp.Name = options.Name
	cliApp.Usage = options.Usage
	cliApp.Version = options.BuildVersion
	cliApp.Authors = options.Users

	if options.Logger == nil {
		if options.Debug == false { // Debug = false
			options.Logger = zap.New(zap.Production)
		} else if options.Quiet { // Debug = true && Quiet = true
			options.Logger = nop.New()
		} else { // Debug = true && Quiet = false
			options.Logger = zap.New(zap.Development)
		}
	}

	if len(options.flags) > 0 {
		cliApp.Flags = append(cliApp.Flags, options.flags...)
	}

	appCommands(options)
	dbCommands(options)

	if len(options.commands) > 0 {
		cliApp.Commands = append(cliApp.Commands, options.commands...)
	}

	sort.Sort(cli.CommandsByName(cliApp.Commands))

	return cliApp.Run(os.Args)
}
