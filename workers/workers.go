package workers

import (
	"context"
	"sync"
	"time"

	"github.com/cryptopay-dev/yaga/graceful"
)

type Workers interface {
	AddJob(time.Duration, func(context.Context))
	AddJobWithError(time.Duration, func(context.Context) error)
	Stop()
	Wait(context.Context) error
}

type workers struct {
	g graceful.Graceful
}

func New() Workers {
	return &workers{
		g: graceful.New(context.Background()),
	}
}

func (w *workers) AddJob(d time.Duration, h func(context.Context)) {
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

// TODO need?
func (w *workers) AddJobWithError(d time.Duration, h func(context.Context) error) {
	w.g.Go(func(ctx context.Context) error {
		tick := time.NewTicker(d)
		defer tick.Stop()

		errCh := make(chan error, 1)
		wg := new(sync.WaitGroup)
		for {
			select {
			case <-tick.C:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := h(ctx); err != nil {
						select {
						case errCh <- err:
						default:
						}
					}
				}()
			case <-ctx.Done():
				wg.Wait()
				return nil
			case err := <-errCh:
				wg.Wait()
				return err
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
