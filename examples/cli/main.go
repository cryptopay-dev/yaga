package main

import (
	"context"
	"os"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/logger/nop"
)

// App structure
type App struct {
}

// Run application
func (App) Run() error { return nil }

// Shutdown application
func (App) Shutdown(ctx context.Context) error { return nil }

func main() {
	instance := App{}

	if cliApp, err := cli.New(
		cli.App(instance),
		cli.Logger(nop.New()),
	); err != nil {
		panic(err)
	} else if err = cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
