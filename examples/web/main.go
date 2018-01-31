package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	e := web.New(web.Options{
		Logger: nop.New(),
		Error: errors.Logic{
			Opts: errors.Options{
				Logger: nop.New(),
			},
		},
		Debug: true,
	})

	e.GET("/test/:command", func(c web.Context) error {
		cmd := c.Param("command")
		fmt.Println("Received command:", cmd)

		switch cmd {
		case "nop":
			// do nothing
		case "off":
			// send signal to exit
			ch <- syscall.SIGABRT
		default:
			// unknown operation
			return http.ErrNotSupported
		}

		return c.JSON(http.StatusOK, cmd)
	})

	go func() {
		if err := web.StartServer(e, os.Getenv("BIND")); err != nil {
			if !strings.Contains(err.Error(), http.ErrServerClosed.Error()) {
				e.Logger.Error(err)
			}
			ch <- syscall.SIGABRT
		}
	}()

	// wait for signals
	sig := <-ch
	fmt.Println("Received signal:", sig.String())

	cancel()

	e.Shutdown(ctx)
}
