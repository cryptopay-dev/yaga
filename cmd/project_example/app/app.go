package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cryptopay-dev/yaga/cli"
	cliApp "github.com/cryptopay-dev/yaga/cli"
	"github.com/cryptopay-dev/yaga/cmd/project_example/app/controllers"
	"github.com/cryptopay-dev/yaga/cmd/project_example/app/library/config"
	"github.com/cryptopay-dev/yaga/errors"
	"github.com/cryptopay-dev/yaga/validate"
	"github.com/cryptopay-dev/yaga/web"
	"gopkg.in/go-playground/validator.v9"
)

// authors scructure
type authors struct {
	Name  string
	Email string
}

// App instance
type App struct {
	cliApp.RunOptions
	Config     config.Config
	LogicError *errors.Logic
	Engine     *web.Engine
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
	return &App{}
}

// Shutdown of application
func (a *App) Shutdown(ctx context.Context) error {
	if a.Engine == nil {
		return nil
	}

	return a.Engine.Shutdown(ctx)
}

// Run of application
func (a *App) Run(opts cliApp.RunOptions) error {
	var err error

	a.RunOptions = opts

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)

	go func() {
		if err = web.StartServer(a.Engine, a.Config.Bind); err != nil {
			a.Logger.Error(err)
			ch <- syscall.SIGABRT
		}
	}()

	// Wait for signals:
	sig := <-ch
	a.Logger.Infof("Received signal: %s", sig.String())

	return a.Shutdown(ctx)
}
