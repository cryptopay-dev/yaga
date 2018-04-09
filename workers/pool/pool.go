package pool

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrEmptyHandler when handler is empty
	ErrEmptyHandler = errors.New("handler must be not null")

	// ErrPoolBusy when pool is busy
	ErrPoolBusy = errors.New("pool is busy")
)

type Pool struct {
	done  chan struct{}
	jobCh chan func(context.Context)
}

func New(size int) *Pool {
	return &Pool{
		done:  make(chan struct{}),
		jobCh: make(chan func(context.Context), size),
	}
}

func (p *Pool) Run(ctx context.Context) {
	size := cap(p.jobCh)
	wg := new(sync.WaitGroup)
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case job := <-p.jobCh:
					job(ctx)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()
	close(p.done)
}

func (p *Pool) Do(job func(context.Context)) error {
	if job == nil {
		return ErrEmptyHandler
	}
	select {
	case p.jobCh <- job:
		return nil
	default:
		return ErrPoolBusy
	}
}

func (p *Pool) Wait(ctx context.Context) error {
	select {
	case <-p.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
