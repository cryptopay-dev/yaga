package main

import (
	"context"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/logger/nop"
)

// App structure
type App struct {
	Config ExampleConfig
}

// Run application
func (App) Run(opts cli.RunOptions) error { return nil }

// Shutdown application
func (App) Shutdown(ctx context.Context) error { return nil }

// ExampleConfig for demo
type ExampleConfig struct {
	Database config.Database `yaml:"database" validate:"required,dive"`
	Redis    config.Redis    `yaml:"redis" validate:"required,dive"`
}

func main() {
	instance := App{}

	if err := cli.Run(
		cli.App(&instance),
		cli.ConfigSource("./config.example.yaml"),
		cli.ConfigInterface(&instance.Config),
		cli.Logger(nop.New()),
	); err != nil {
		panic(err)
	}
}
