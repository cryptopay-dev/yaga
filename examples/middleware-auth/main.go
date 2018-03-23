package main

import (
	"net/http"

	"github.com/cryptopay-dev/yaga/helpers/postgres"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/middlewares/auth"
	"github.com/cryptopay-dev/yaga/web"
)

func main() {
	log := nop.New()
	e, err := web.New(web.Options{
		Logger: log,
	})

	if err != nil {
		log.Panic(err)
	}

	db, err := postgres.Connect("database")
	if err != nil {
		log.Panic(err)
	}

	authenticate := auth.New(
		auth.Logger(log),
		auth.DB(db),
	)

	e.Use(authenticate.Middleware())

	e.GET("/", func(ctx web.Context) error {
		return ctx.String(
			http.StatusOK,
			"Private zone",
		)
	}, authenticate.Middleware())

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}
