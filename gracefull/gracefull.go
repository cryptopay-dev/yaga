package gracefull

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/sync/errgroup"
)

type logger interface {
	Infof(string, ...interface{})
}

func GracefullShutdown(ctx context.Context, log logger) (*errgroup.Group, context.Context) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)

	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer cancel()
		sig := <-ch
		if log != nil {
			log.Infof("received signal: %s", sig.String())
		}
		return nil
	})

	return g, ctx
}

func Gracefull(ctx context.Context) (*errgroup.Group, context.Context) {
	return errgroup.WithContext(ctx)
}
