package main

import (
	"net/http"
	"time"

	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/report"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	e := web.New(web.Options{
		Logger: nop.New(),
		Error: errors.Logic{Opts: errors.Options{
			Logger: nop.New(),
		}},
	})

	e.GET("/", func(ctx web.Context) error {
		ts := time.Now()

		// Variant 1:
		defer report.CaptureResponseTimings(
			ctx.Response(),
			ctx.Path(),
			time.Now(),
		)

		// Variant 2:
		defer report.CaptureResponseTime(
			"some-platform",
			"some-action",
			time.Now(),
		)

		// Capture margin:
		defer report.CaptureMargin(
			"some-platform",
			"some-pair",
			float64(time.Since(ts)),
		)

		return ctx.String(
			http.StatusOK,
			"OK",
		)
	})

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
