package app

import (
	"context"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/cmd/yaga/project_example/app/controllers"
	"github.com/cryptopay-dev/yaga/cmd/yaga/project_example/app/library/config"
	"github.com/cryptopay-dev/yaga/graceful"
	"github.com/cryptopay-dev/yaga/web"
)

// authors scructure
type authors struct {
	Name  string
	Email string
}

// App instance
type App struct {
	cli.RunOptions
	Config   config.Config
	Engine   *web.Engine
	Graceful graceful.Graceful
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
		Graceful: graceful.New(context.Background()),
	}
}

// Shutdown of application
func (a *App) Shutdown(ctx context.Context) error {
	if a.Engine == nil {
		return nil
	}

	return a.Engine.Shutdown(ctx)
}

// Run of application
func (a *App) Run(opts cli.RunOptions) error {
	var err error

	a.RunOptions = opts

	if a.Engine, err = web.New(web.Options{
		Logger: a.Logger,
		Debug:  a.Debug,
	}); err != nil {
		return err
	}

	if _, err = controllers.New(
		controllers.Logger(a.Logger),
		controllers.Config(&a.Config),
		controllers.Engine(a.Engine),
		controllers.BuildTime(a.BuildTime),
		controllers.BuildVersion(a.BuildVersion),
	); err != nil {
		return err
	}

	graceful.AttachNotifier(a.Graceful, a.Logger)
	web.StartAsync(a.Engine, a.Config.Bind, a.Graceful)

	return a.Graceful.Wait()
}
