package workers

import (
	"sync"
	"time"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/robfig/cron"
)

// The Schedule describes a job's duty cycle.
type Schedule = cron.Schedule

// ConstantDelaySchedule represents a simple recurring duty cycle, e.g. "Every 5 minutes".
type ConstantDelaySchedule = cron.ConstantDelaySchedule

// Cmd handler to run
type Cmd = func() error

// Entry consists of a schedule and the func to execute on that schedule.
type Entry = cron.Entry

// Workers struct
type Workers struct {
	wg   sync.WaitGroup
	cron *cron.Cron
}

// New returns a new workers runner.
func New() *Workers {
	w := &Workers{
		wg:   sync.WaitGroup{},
		cron: cron.New(),
	}

	return w
}

// AddFunc adds a func to the Cron to be run on the given schedule.
func (w *Workers) AddFunc(spec string, cmd Cmd) error {
	return w.cron.AddFunc(spec, func() {
		w.wg.Add(1)
		if err := cmd(); err != nil {
			log.Error(err)
		}
		w.wg.Done()
	})
}

// Start the cron scheduler in its own go-routine.
func (w *Workers) Start() { w.cron.Start() }

// Stop the cron scheduler and wait for jobs.
func (w *Workers) Stop() {
	w.cron.Stop()
	w.wg.Wait()
}

// Entries returns a snapshot of the cron entries.
func (w *Workers) Entries() []*Entry { return w.cron.Entries() }

// Every returns a crontab Schedule that activates once every duration.
func Every(duration time.Duration) ConstantDelaySchedule { return cron.Every(duration) }
