package workers

import (
	"context"
	"sync"
	"time"

	"github.com/cryptopay-dev/yaga/graceful"
)

type Handler func(context.Context)

type workers struct {
	g graceful.Graceful
}

func New() *workers {
	return &workers{
		g: graceful.New(context.Background()),
	}
}

func (w *workers) AddJob(d time.Duration, h Handler) {
	w.g.Go(func(ctx context.Context) error {
		tick := time.NewTicker(d)
		defer tick.Stop()

		wg := new(sync.WaitGroup)
		for {
			select {
			case <-tick.C:
				wg.Add(1)
				go func() {
					defer wg.Done()
					h(ctx)
				}()
			case <-ctx.Done():
				wg.Wait()
				return nil
			}
		}
	})
}

func (w *workers) Stop() {
	w.g.Cancel()
}

func (w *workers) Wait(ctx context.Context) error {
	return w.g.Wait(ctx)
}
