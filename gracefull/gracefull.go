package gracefull

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/sync/errgroup"
)

type Gracefull interface {
	Go(func() error)
	Wait() error
	Cancel()
}

type logger interface {
	Infof(string, ...interface{})
}

type gracefull struct {
	*errgroup.Group
	cancel context.CancelFunc
}

func (g *gracefull) Cancel() {
	g.cancel()
}

func NewNotify(ctx context.Context, log logger) (Gracefull, context.Context) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)

	g, ctx := New(ctx)
	g.Go(func() error {
		defer g.Cancel()
		sig := <-ch
		if log != nil {
			log.Infof("received signal: %s", sig.String())
		}
		return nil
	})

	return g, ctx
}

func New(ctx context.Context) (Gracefull, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	return &gracefull{g, cancel}, ctx
}
