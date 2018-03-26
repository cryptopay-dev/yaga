package main

import (
	"context"
	"net/http"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/graceful"
	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	log := zap.New(zap.Development)

	e, err := web.New(web.Options{
		Logger: log,
		Debug:  true,
	})

	if err != nil {
		log.Panic(err)
		return
	}

	e.GET("/test/:command", func(c web.Context) error {
		cmd := c.Param("command")
		log.Infof("Received command: %v", cmd)

		switch cmd {
		case "nop":
			// do nothing
		default:
			// unknown operation
			return http.ErrNotSupported
		}

		return c.JSON(http.StatusOK, cmd)
	})

	g := graceful.New(context.Background())
	graceful.AttachNotifier(g, e.Logger)

	web.StartAsync(e, config.GetString("bind"), g)

	if err := g.Wait(); err != nil {
		e.Logger.Error(err)
	}
}
