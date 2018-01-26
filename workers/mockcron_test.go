package workers

import (
	"sync"
	"time"
)

type dummySchedule struct{}

func (dummySchedule) Next(t time.Time) time.Time {
	return t.Add(time.Minute)
}

type mockCron struct {
	mu      sync.Mutex
	tickers []*time.Ticker
	stopped bool

	poolWorker *pool
	stop       chan struct{}
}

func (m *mockCron) New(duration time.Duration, job func()) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-m.stop:
				return
			case <-ticker.C:
				go job()
			}
		}
	}()
	m.tickers = append(m.tickers, ticker)
}

func (m *mockCron) Start() {
	m.poolWorker.stop.Store(false)
}

func (m *mockCron) Stop() {
	m.poolWorker.stop.Store(true)
}

func (m *mockCron) StopCron() {
	m.poolWorker.stop.Store(true)

	m.mu.Lock()
	defer m.mu.Unlock()
	if m.stopped {
		return
	}

	close(m.stop)
	for _, ticker := range m.tickers {
		ticker.Stop()
	}
	m.tickers = nil
	m.stopped = true
}

func (m *mockCron) Wait() {
	m.poolWorker.wg.Wait()
}

func newCronForTest() (*mockCron, func(string, time.Duration, func()) (*worker, error)) {
	c := &mockCron{
		poolWorker: newPool(),
		stop:       make(chan struct{}),
	}

	return c, func(name string, tick time.Duration, handler func()) (*worker, error) {
		opts := Options{
			Name:     name,
			Schedule: dummySchedule{},
			Handler:  handler,
		}

		return newWorker(opts, c.poolWorker, func(w *worker, f func()) {
			c.New(tick, f)
		})
	}
}
