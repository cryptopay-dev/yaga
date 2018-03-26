package main

import (
	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/cmd/yaga/project_example/app"
	"github.com/cryptopay-dev/yaga/cmd/yaga/project_example/misc"
)

func main() {
	instance := app.New()

	if err := cli.Run(
		cli.App(instance),
		cli.Debug(misc.Debug),
		cli.Name(misc.Name),
		cli.Usage(misc.Usage),
		cli.Users(app.Authors()),
		cli.BuildTime(misc.BuildTime),
		cli.BuildVersion(misc.Version),
	); err != nil {
		panic(err)
	}
}
