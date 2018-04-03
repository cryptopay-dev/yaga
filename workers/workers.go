package workers

import (
	"sync"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/robfig/cron"
)

// Cmd handler to run
type Cmd = func() error

// Entry consists of a schedule and the func to execute on that schedule.
type Entry = cron.Entry

// Workers struct
type Workers struct {
	wg   *sync.WaitGroup
	cron *cron.Cron
}

// New returns a new workers runner.
func New() *Workers {
	w := &Workers{
		wg:   new(sync.WaitGroup),
		cron: cron.New(),
	}

	return w
}

// AddFunc adds a func to the Cron to be run on the given schedule.
func (w *Workers) AddFunc(spec string, cmd Cmd) error {
	return w.cron.AddFunc(spec, func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("workers panic: %v", r)
			}
		}()

		w.wg.Add(1)
		defer w.wg.Done()

		if err := cmd(); err != nil {
			log.Error(err)
		}
	})
}

// Start the cron scheduler in its own go-routine.
func (w *Workers) Start() { w.cron.Start() }

// Stop the cron scheduler and wait for jobs.
func (w *Workers) Stop() {
	w.cron.Stop()
	w.wg.Wait()
}
