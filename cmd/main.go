package main

import (
	"os"
	"strings"
	"time"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/config"
)

const (
	Name  = "yaga"
	Usage = "initialized service"
)

var (
	// Version of application by default is 'dev'
	Version = "dev"
	// BuildTime of application by default - current time
	BuildTime = time.Now().Format(time.RFC3339)

	Config struct {
		Database config.Database `yaml:"database" validate:"required"`
	}
)

func beforeHandler() cli.Handler {
	return func(opts *cli.Options) func(c *cli.Context) error {
		return func(c *cli.Context) error {
			var (
				configSrc interface{}
				configDst interface{}

				filename = c.String("config")
			)

			if len(filename) > 0 {
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					// does not exist
					if c.IsSet("config") {
						return err
					}
				} else {
					// file exists
					configSrc = filename
					configDst = &Config
				}
			}

			if configSrc == nil {
				configSrc = strings.NewReader("")
				configDst = &struct{}{}
			}

			cli.Config(configSrc, configDst)(opts)

			return nil
		}
	}
}

func projectHandler() cli.Commandor {
	setWorkDir := cli.StringFlag{
		Name:  "dir",
		Value: "",
		Usage: "workdir and name project",
	}

	return func(opts *cli.Options) (c cli.Command) {
		c.Name = "init"
		c.Aliases = []string{"i"}
		c.Usage = "Creating new project"
		c.Flags = []cli.Flag{setWorkDir}
		c.Action = func(ctx *cli.Context) (err error) {
			return projectBuilder(opts, ctx.String(setWorkDir.Name))
		}

		return
	}
}

func main() {
	if err := cli.Run(
		cli.Name(Name),
		cli.Usage(Usage),
		cli.Debug(true, false),
		cli.BuildTime(BuildTime),
		cli.BuildVersion(Version),
		cli.Flags(func(*cli.Options) cli.Flag {
			return cli.StringFlag{
				Name:  "config",
				Usage: "config for migrations",
				Value: "./config.yaml",
			}
		}),
		cli.Commands(projectHandler()),
		cli.Trigger(nil, beforeHandler(), nil),
	); err != nil {
		panic(err)
	}
}
