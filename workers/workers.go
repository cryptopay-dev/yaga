package workers

import (
	"errors"
	"time"

	"github.com/robfig/cron"
)

type (
	// Options structure for creation new worker.
	Options struct {
		Name     string
		Schedule Schedule
		Handler  func()
	}

	// Schedule describes a job's duty cycle.
	//
	// Return the next activation time, later than the given time.
	// Next is invoked initially, and then each time the job is run.
	Schedule = cron.Schedule
)

var (
	// ErrAlreadyWorker is returned by New calls
	// when workers name is already exists.
	ErrAlreadyWorker = errors.New("worker name must be unique")

	// ErrWrongOptions is returned by New calls
	// when parameter Options.Schedule is NIL or Options.Handler is NIL.
	ErrWrongOptions = errors.New("wrong options")

	cronWorker = cron.New()
	poolWorker = newPool()
)

// New returns an error if cannot create new worker
func New(opts Options) (err error) {
	_, err = newWorker(opts, poolWorker, func(w *worker, handler func()) {
		cronWorker.Schedule(w.options.Schedule, cron.FuncJob(handler))
	})

	return
}

// Parse returns a new crontab schedule representing the given spec.
// It returns a descriptive error if the spec is not valid.
//
// It accepts
//   - Full crontab specs, e.g. "* * * * * ?"
//   - Descriptors, e.g. "@midnight", "@every 1h30m"
func Parse(spec string) (Schedule, error) {
	return cron.Parse(spec)
}

// Every returns a crontab Schedule that activates once every duration.
// Delays of less than a second are not supported (will round up to 1 second).
// Any fields less than a Second are truncated.
func Every(duration time.Duration) Schedule {
	return cron.Every(duration)
}

// Start all workers.
func Start() {
	poolWorker.stop.Store(false)
	cronWorker.Start()
}

// Stop all workers.
func Stop() {
	poolWorker.stop.Store(true)
	cronWorker.Stop()
}

// Wait blocks until all workers will be stopped.
func Wait() {
	poolWorker.wg.Wait()
}
