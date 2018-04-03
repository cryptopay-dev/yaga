package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/sync/errgroup"
	"github.com/pkg/errors"
)

// Graceful interface
type Graceful interface {
	Go(func(context.Context) error)
	Wait(context.Context) error
	Context() context.Context
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

func (g *graceful) Context() context.Context {
	return g.ctx
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
			err = errors.Wrap(err, "graceful failed")
		}()
		err = job(g.ctx)
		return
	}

	g.eg.Go(f)
}

func (g *graceful) Wait(ctx context.Context) error {
	if ctx == nil {
		return g.eg.Wait()
	}

	var err error
	done := make(chan struct{})
	go func() {
		err = g.eg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	g.Go(func(c context.Context) error {
		select {
		case sig := <-ch:
			defer g.Cancel()
			if log != nil {
				log.Infof("received signal: %s", sig.String())
			}
		case <-c.Done():
			return c.Err()
		}
		return nil
	})
}
