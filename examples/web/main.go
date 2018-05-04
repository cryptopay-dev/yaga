package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/graceful"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	e, err := web.New(web.Options{Debug: true})

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

	web.StartAsync(e)

	graceful.Wait()
}
