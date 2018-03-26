package main

import (
	"context"
	"fmt"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/config"
	"github.com/davecgh/go-spew/spew"
)

// App structure
type App struct{}

// Run application
func (App) Run(opts cli.RunOptions) error {
	spew.Dump(opts)
	return nil
}

// Shutdown application
func (App) Shutdown(ctx context.Context) error { return nil }

var echo string

func main() {
	instance := App{}
	config.SetConfigName("config.example")

	if err := cli.Run(
		cli.App(&instance),
		cli.Debug(true, true), // Debug & Quiet
		cli.Flags(cli.StringFlag{
			Name:        "echo",
			Usage:       "echo printing",
			Destination: &echo,
		}),
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
