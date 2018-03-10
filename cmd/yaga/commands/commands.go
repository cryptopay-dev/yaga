package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/labstack/gommon/color"
	"github.com/urfave/cli"
	"go.uber.org/atomic"
)

var (
	clr = color.New()
	cnt = atomic.NewInt64(0)
)

// All returns all commands
func All() cli.Commands {
	clr.Enable()

	return []cli.Command{
		newProject(), // Creates new project..
	}
}

type formatter = func(msg interface{}, styles ...string) string

func output(out io.Writer, format formatter, msg string) {
	clr.SetOutput(out)
	cnt.Add(1)
	clr.Printf("[%03d] %s\n", cnt.Load(), format(msg))
}

func print(msg ...interface{}) {
	output(os.Stdout, clr.Blue, fmt.Sprint(msg...))
}

func printf(format string, msg ...interface{}) {
	output(os.Stdout, clr.Blue, fmt.Sprintf(format, msg...))
}

func info(msg ...interface{}) {
	output(os.Stdout, clr.Green, fmt.Sprint(msg...))
}

func infof(format string, msg ...interface{}) {
	output(os.Stdout, clr.Green, fmt.Sprintf(format, msg...))
}

func errors(msg ...interface{}) {
	output(os.Stderr, clr.Red, fmt.Sprint(msg...))
	os.Exit(1)
}

func errorsf(format string, msg ...interface{}) {
	output(os.Stderr, clr.Red, fmt.Sprintf(format, msg...))
	os.Exit(1)
}
