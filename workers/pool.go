package workers

import (
	"sync"

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
	jobCh  chan func()
	waitCh chan chan struct{}
}

func newPool() *pool {
	p := &pool{
		workers: make(map[string]*worker),
		running: atomic.NewBool(false),

		cmdCh:  make(chan commandAction, 1),
		jobCh:  make(chan func(), 8),
		waitCh: make(chan chan struct{}),
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

		select {
		case p.jobCh <- w.options.Handler:
		default:
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

func (p *pool) wait() {
	<-<-p.waitCh
}

func (p *pool) dispatcher() {
	var (
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
			go func() {
				defer wg.Done()
				job()
			}()
		case p.waitCh <- waiter:
		case cmd := <-p.cmdCh:
			switch cmd {
			case stop:
				wg.Wait()
				close(waiter)
				running = false
			case start:
				waiter = make(chan struct{})
				running = true
			default:
			}
		}
	}
}
