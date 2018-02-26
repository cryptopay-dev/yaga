package main

import (
	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/cmd/project_example/app"
	"github.com/cryptopay-dev/yaga/cmd/project_example/vars"
)

func main() {
	instance := app.New()

	if err := cli.Run(
		cli.App(instance),
		cli.Config(vars.Config, &instance.Config),
		cli.Debug(vars.Debug),
		cli.Name(vars.Name),
		cli.Usage(vars.Usage),
		cli.Users(app.Authors()),
		cli.BuildTime(vars.BuildTime),
		cli.BuildVersion(vars.Version),
	); err != nil {
		panic(err)
	}
}
