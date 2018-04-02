package workers

import (
	"context"
	"sync"
	stdAtomic "sync/atomic"

	"go.uber.org/atomic"
)

type commandAction int

const (
	stop commandAction = iota
	start
)

type pool struct {
	running *atomic.Bool

	mu      sync.Mutex
	workers map[string]*worker

	cmdCh  chan commandAction
	jobCh  chan func(context.Context)
	waiter stdAtomic.Value
}

func newPool() *pool {
	p := &pool{
		workers: make(map[string]*worker),
		running: atomic.NewBool(false),

		cmdCh: make(chan commandAction, 1),
		jobCh: make(chan func(context.Context), 8),
	}
	go p.dispatcher()

	return p
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
		if !p.running.Load() {
			return
		}

		for {
			select {
			case p.jobCh <- w.options.Handler:
			default:
				if !p.running.Load() {
					return
				}
			}
		}
	}

	return w, nil
}

func (p *pool) start() {
	if !p.running.Swap(true) {
		p.cmdCh <- start
	}
}

func (p *pool) stop() {
	if p.running.Swap(false) {
		p.cmdCh <- stop
	}
}

func (p *pool) wait(ctx context.Context) error {
	waiter, ok := p.waiter.Load().(chan struct{})
	if !ok {
		return nil
	}
	if ctx == nil {
		<-waiter
		return nil
	}
	select {
	case <-waiter:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *pool) dispatcher() {
	var (
		ctx     context.Context
		cancel  context.CancelFunc
		running bool
		waiter  chan struct{}
		wg      = new(sync.WaitGroup)
	)

	for {
		select {
		case job := <-p.jobCh:
			if !running {
				continue
			}
			wg.Add(1)
			go func(c context.Context) {
				defer wg.Done()
				job(c)
			}(ctx)
		case cmd := <-p.cmdCh:
			switch cmd {
			case stop:
				cancel()
				wg.Wait()
				close(waiter)
				running = false
			case start:
				waiter = make(chan struct{})
				p.waiter.Store(waiter)
				ctx, cancel = context.WithCancel(context.Background())
				running = true
			default:
			}
		}
	}
}
