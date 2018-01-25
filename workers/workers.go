package workers

import (
	"errors"
	"time"

	"github.com/robfig/cron"
)

type (
	CronHandler func(Schedule, func())

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

func newWorker(opts Options, p *pool, addToCron CronHandler) (*worker, error) {
	if opts.Schedule == nil || opts.Handler == nil {
		return nil, ErrWrongOptions
	}

	w, err := p.createWorker(opts)
	if err != nil {
		return nil, err
	}
	w.pool = p

	addToCron(w.options.Schedule, w.job)

	return w, nil
}

func New(opts Options) (err error) {
	_, err = newWorker(opts, poolWorker, func(schedule Schedule, handler func()) {
		cronWorker.Schedule(schedule, cron.FuncJob(handler))
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
