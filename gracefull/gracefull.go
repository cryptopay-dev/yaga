package gracefull

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/sync/errgroup"
)

// Gracefull interface
type Gracefull interface {
	Go(func(context.Context) error)
	Wait() error
	Cancel()
}

type logger interface {
	Infof(string, ...interface{})
}

type gracefull struct {
	eg     *errgroup.Group
	cancel context.CancelFunc
	ctx    context.Context
}

func (g *gracefull) Cancel() {
	g.cancel()
}

func (g *gracefull) Go(job func(context.Context) error) {
	f := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("Unknown panic")
				}
			}
		}()

		err = job(g.ctx)
		return
	}

	g.eg.Go(f)
}

func (g *gracefull) Wait() error {
	return g.eg.Wait()
}

// NewNotify returns a new Gracefull and an associated Context derived from ctx.
//
// Returns Gracefull associated with notification of OS signals
func NewNotify(ctx context.Context, log logger) Gracefull {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)

	g := New(ctx)
	g.Go(func(c context.Context) error {
		select {
		case sig := <-ch:
			defer g.Cancel()
			if log != nil {
				log.Infof("received signal: %s", sig.String())
			}
		case <-c.Done():
		}
		return nil
	})

	return g
}

// New returns a new Gracefull and an associated Context derived from ctx.
func New(ctx context.Context) Gracefull {
	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	return &gracefull{g, cancel, ctx}
}
