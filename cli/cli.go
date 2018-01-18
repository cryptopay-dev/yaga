package cli

import (
	"sort"

	"github.com/urfave/cli"
)

func New(opts ...Option) (*cli.App, error) {
	var options, err = newOptions(opts...)

	if err != nil {
		return nil, err
	}

	cliApp := cli.NewApp()
	cliApp.Name = options.Name
	cliApp.Usage = options.Usage
	cliApp.Version = options.BuildVersion
	cliApp.Authors = options.Users

	addCommands(cliApp, options)
	sort.Sort(cli.CommandsByName(cliApp.Commands))
	return cliApp, nil
}
