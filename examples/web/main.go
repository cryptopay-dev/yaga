package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := nop.New()

	e, err := web.New(web.Options{
		Logger: log,
		Debug:  true,
	})

	if err != nil {
		log.Panic(err)
	}

	e.GET("/test/:command", func(c web.Context) error {
		cmd := c.Param("command")
		fmt.Println("Received command:", cmd)

		switch cmd {
		case "nop":
			// do nothing
		default:
			// unknown operation
			return http.ErrNotSupported
		}

		return c.JSON(http.StatusOK, cmd)
	})

	done := web.StartAsync(e, config.GetString("bind"))

	// wait for signals
	sig := <-done
	log.Info("Received signal:", sig.String())

	e.Shutdown(ctx)
}
