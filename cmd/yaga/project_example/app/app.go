package app

import (
	"context"
	"time"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/cmd/yaga/project_example/app/controllers"
	"github.com/cryptopay-dev/yaga/cmd/yaga/project_example/app/library/config"
	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/graceful"
	"github.com/cryptopay-dev/yaga/helpers/shutdown"
	"github.com/cryptopay-dev/yaga/validate"
	"github.com/cryptopay-dev/yaga/web"
	"github.com/cryptopay-dev/yaga/workers"
	"gopkg.in/go-playground/validator.v9"
)

// authors scructure
type authors struct {
	Name  string
	Email string
}

// App instance
type App struct {
	cli.RunOptions
	Config     config.Config
	LogicError *errors.Logic
	Engine     *web.Engine
	Graceful   graceful.Graceful
	Workers    workers.Workers
}

var appAuthors = []authors{
	{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	},
}

// Authors of application
func Authors() []cli.Author {
	var result = make([]cli.Author, 0, len(appAuthors))
	for i := range appAuthors {
		result = append(result, cli.Author(appAuthors[i]))
	}

	return result
}

// New creates instance
func New() *App {
	return &App{
		Workers: workers.New(),
	}
}

// Shutdown of application
func (a *App) Shutdown(ctx context.Context) error {
	return a.Graceful.Wait(ctx)
}

// Run of application
func (a *App) Run(opts cli.RunOptions) error {
	var err error

	a.RunOptions = opts

	ctx := shutdown.ShutdownContext(context.Background(), a.Logger)
	a.Graceful = graceful.New(ctx)

	if a.LogicError, err = errors.New(errors.Options{
		Debug:  a.Debug,
		Logger: a.Logger,
	}); err != nil {
		return err
	}

	v := validator.New()

	a.Engine = web.New(web.Options{
		Logger:    a.Logger,
		Error:     a.LogicError,
		Debug:     a.Debug,
		Validator: validate.New(v),
	})

	controllers.New(
		controllers.Logger(a.Logger),
		controllers.Config(&a.Config),
		controllers.Engine(a.Engine),
		controllers.BuildTime(a.BuildTime),
		controllers.BuildVersion(a.BuildVersion),
	)

	web.StartAsync(a.Engine, a.Config.Bind, a.Graceful)
	workers.AttachGraceful(a.Workers, a.Graceful, time.Second*30)

	<-ctx.Done()

	// TODO return a.Shutdown()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return a.Shutdown(ctx)
}
