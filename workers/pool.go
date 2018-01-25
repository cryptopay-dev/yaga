package workers

import (
	"sync"

	"go.uber.org/atomic"
)

type pool struct {
	workers map[string]*worker
	mu      sync.Locker
	stop    *atomic.Bool
	wg      sync.WaitGroup
}

func newPool() *pool {
	return &pool{
		workers: make(map[string]*worker),
		mu:      new(sync.Mutex),
		stop:    atomic.NewBool(false),
		wg:      sync.WaitGroup{},
	}
}

func (p *pool) createWorker(opts Options) (*worker, error) {
	w := &worker{
		options: opts,
	}

	p.mu.Lock()
	if _, found := p.workers[opts.Name]; found {
		p.mu.Unlock()
		return nil, ErrAlreadyWorker
	}
	p.workers[opts.Name] = w
	p.mu.Unlock()

	w.job = func() {
		if p.stop.Load() {
			return
		}
		p.wg.Add(1)
		defer p.wg.Done()

		w.options.Handler()
	}

	return w, nil
}
