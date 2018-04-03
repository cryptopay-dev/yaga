package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/helpers/shutdown"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	log := nop.New()
	ctx := shutdown.ShutdownContext(context.Background(), log)

	errLogic, _ := errors.New(errors.Options{
		Logger: log,
	})

	e := web.New(web.Options{
		Logger: log,
		Error:  errLogic,
		Debug:  true,
	})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	e.GET("/test/:command", func(c web.Context) error {
		cmd := c.Param("command")
		fmt.Println("Received command:", cmd)

		switch cmd {
		case "nop":
			// do nothing
		case "off":
			// send signal to exit
			cancel()
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
			cancel()
		}
	}()

	// wait for signals
	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	e.Shutdown(ctx)
}
