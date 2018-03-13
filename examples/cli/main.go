package main

import (
	"context"
	"fmt"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/config"
	"github.com/davecgh/go-spew/spew"
)

// App structure
type App struct {
	Config ExampleConfig
}

// Run application
func (App) Run(opts cli.RunOptions) error {
	spew.Dump(opts.DB)
	spew.Dump(opts.Redis)
	return nil
}

// Shutdown application
func (App) Shutdown(ctx context.Context) error { return nil }

// ExampleConfig for demo
type ExampleConfig struct {
	Database config.Database `yaml:"database" validate:"required,dive"`
	Redis    config.Redis    `yaml:"redis" validate:"required,dive"`
}

// Config path to config.example.yaml
const Config = "./examples/cli/config.example.yaml"

var echo string

func beforeHandler() cli.Handler {
	return func(opts *cli.Options) func(c *cli.Context) error {
		return func(c *cli.Context) error {
			val := c.String("echo")
			if len(val) > 0 {
				fmt.Println("print echo:", val)
			}

			return nil
		}
	}
}

func main() {
	instance := App{}

	if err := cli.Run(
		cli.App(&instance),
		cli.Config(Config, &instance.Config),
		cli.Debug(true, true), // Debug & Quiet
		cli.Flags(cli.StringFlag{
			Name:        "echo",
			Usage:       "echo printing",
			Destination: &echo,
		}),
		cli.Trigger(nil, beforeHandler(), nil),
		cli.Commands(func(opts *cli.Options) (c cli.Command) {
			c.Name = "test"
			c.Aliases = []string{"t"}
			c.Usage = "run test command"
			c.Action = func(c *cli.Context) error {
				if len(echo) > 0 {
					fmt.Println("echo:", echo)
				}
				fmt.Println("test: Hello world!")
				return nil
			}
			return
		}),
	); err != nil {
		panic(err)
	}
}
