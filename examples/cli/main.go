package main

import (
	"context"
	"fmt"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/config"
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

// Config path to config.example.yaml
const Config = "./config.example.yaml"

func main() {
	instance := App{}

	if err := cli.Run(
		cli.App(&instance),
		cli.Config(Config, &instance.Config),
		cli.Debug(true, true), // Debug & Quiet
		cli.Commands(func(opts *cli.Options) (c cli.Command) {
			c.Name = "test"
			c.Aliases = []string{"t"}
			c.Usage = "run test command"
			c.Action = func(c *cli.Context) error {
				fmt.Println("test: Hello world!")
				return nil
			}
			return
		}),
	); err != nil {
		panic(err)
	}
}
