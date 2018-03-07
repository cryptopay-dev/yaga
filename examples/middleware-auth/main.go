package main

import (
	"net/http"
	"os"

	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/middlewares/auth"
	"github.com/cryptopay-dev/yaga/web"
	"github.com/go-pg/pg"
)

func main() {
	log := nop.New()
	e, err := web.New(web.Options{
		Logger: log,
	})

	if err != nil {
		log.Panic(err)
	}

	authenticate := auth.New(
		auth.Logger(log),
		auth.DB(pg.Connect(&pg.Options{
			Addr:     os.Getenv("DATABASE_ADDR"),
			User:     os.Getenv("DATABASE_USER"),
			Database: os.Getenv("DATABASE_DATABASE"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			PoolSize: 2,
		})),
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
