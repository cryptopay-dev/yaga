package main

import (
	"net/http"
	"os"

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

	ctx := web.StartAsync(e, os.Getenv("BIND"))

	// wait for signals
	<-ctx.Done()
	log.Info("Stopping...")

	e.Shutdown(ctx)
	log.Info("Shutdown")
}
