package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/sync/errgroup"
)

// Graceful interface
type Graceful interface {
	Go(func(context.Context) error)
	Wait() error
	Cancel()
}

type logger interface {
	Infof(string, ...interface{})
}

type graceful struct {
	eg     *errgroup.Group
	cancel context.CancelFunc
	ctx    context.Context
}

func (g *graceful) Cancel() {
	g.cancel()
}

func (g *graceful) Go(job func(context.Context) error) {
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

func (g *graceful) Wait() error {
	return g.eg.Wait()
}

// New returns a new Graceful and an associated Context derived from ctx.
func New(ctx context.Context) Graceful {
	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	return &graceful{g, cancel, ctx}
}

// AttachNotifier connects Graceful to notification of OS signals.
func AttachNotifier(g Graceful, log logger) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGABRT)

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
}
