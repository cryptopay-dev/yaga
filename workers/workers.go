package workers

import (
	"errors"
	"time"

	"github.com/robfig/cron"
)

type (
	Options struct {
		Name     string
		Schedule Schedule
		Handler  func()
	}

	Schedule interface {
		cron.Schedule
	}
)

var (
	ErrAlreadyWorker = errors.New("worker name must be unique")
	ErrWrongOptions  = errors.New("wrong options")

	cronWorker = cron.New()
	poolWorker = newPool()
)

func New(opts Options) (err error) {
	_, err = newWorker(opts, poolWorker, func(w *worker, handler func()) {
		cronWorker.Schedule(w.options.Schedule, cron.FuncJob(handler))
	})

	return
}

func Parse(spec string) (Schedule, error) {
	return cron.Parse(spec)
}

func Every(duration time.Duration) Schedule {
	return cron.Every(duration)
}

func Start() {
	poolWorker.stop.Store(false)
	cronWorker.Start()
}

func Stop() {
	poolWorker.stop.Store(true)
	cronWorker.Stop()
}

func Wait() {
	poolWorker.wg.Wait()
}
