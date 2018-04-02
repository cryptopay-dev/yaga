package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cryptopay-dev/yaga/cmd/yaga/commands"
	"github.com/cryptopay-dev/yaga/cmd/yaga/internal"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/urfave/cli"
)

const (
	// Name of application
	Name = "yaga"
	// Usage of application
	Usage = "Yaga command line tool"
)

var (
	// Version of application by default is 'dev'
	Version = "dev"
	// BuildTime of application by default - current time
	BuildTime = time.Now().Format(time.RFC3339)
	// format Version and BuildTime
	format = "%s (%s)"
)

func init() {
	ver, dt, err := internal.Version()
	if err != nil {
		return
	}

	Version = ver
	BuildTime = dt
}

func main() {
	log.New()
	a := cli.NewApp()
	a.Name = Name
	a.Usage = Usage
	a.Version = fmt.Sprintf(format, Version, BuildTime)
	a.Commands = commands.All()
	if err := a.Run(os.Args); err != nil {
		panic(err)
	}
}
