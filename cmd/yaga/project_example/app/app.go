package app

import (
	"context"

	"github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/cmd/yaga/project_example/app/controllers"
	"github.com/cryptopay-dev/yaga/graceful"
	"github.com/cryptopay-dev/yaga/web"
)

// App instance
type App struct {
	cli.RunOptions
	Engine *web.Engine
}

// Authors of application
func Authors() []cli.Author {
	return []cli.Author{
		{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
	}
}

// New creates instance
func New() *App { return &App{} }

// Shutdown of application
func (a *App) Shutdown(ctx context.Context) error { return nil }

// Run of application
func (a *App) Run(opts cli.RunOptions) error {
	var err error

	a.RunOptions = opts

	if a.Engine, err = web.New(web.Options{Debug: a.Debug}); err != nil {
		return err
	}

	if _, err = controllers.New(
		controllers.Engine(a.Engine),
		controllers.BuildTime(a.BuildTime),
		controllers.BuildVersion(a.BuildVersion),
	); err != nil {
		return err
	}

	web.StartAsync(a.Engine)

	return graceful.Wait()
}
