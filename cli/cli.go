package cli

import (
	"errors"
	"os"
	"reflect"
	"sort"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/urfave/cli"
	"gopkg.in/go-playground/validator.v9"
)

var (
	ErrAppNotPointer    = errors.New("app must be an pointer to an struct")
	ErrConfigNotPointer = errors.New("config must be an pointer to an struct")
)

// Run creates instance of cli.App with Options.
// Validate options with https://github.com/go-playground/validator
// Required:
// - App instance
// - Logger
func Run(opts ...Option) error {
	var (
		err     error
		options = newOptions(opts...)
	)

	// If we have config-source/interface - loading config:
	if options.ConfigSource != nil &&
		options.ConfigInterface != nil {
		if reflect.TypeOf(options.ConfigInterface).Kind() != reflect.Ptr {
			return ErrConfigNotPointer
		}

		if err = config.Load(
			options.ConfigSource,
			options.ConfigInterface,
		); err != nil {
			return err
		}
	}

	if reflect.TypeOf(options.App).Kind() != reflect.Ptr {
		return ErrAppNotPointer
	}

	// Validate options:
	if err = validator.New().Struct(options); err != nil {
		return err
	}

	cliApp := cli.NewApp()
	cliApp.Name = options.Name
	cliApp.Usage = options.Usage
	cliApp.Version = options.BuildVersion
	cliApp.Authors = options.Users

	addCommands(cliApp, options)
	sort.Sort(cli.CommandsByName(cliApp.Commands))

	return cliApp.Run(os.Args)
}
