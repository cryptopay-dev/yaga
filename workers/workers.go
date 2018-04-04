package workers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/logger/log"
	wrap "github.com/pkg/errors"
	"github.com/robfig/cron"
)

// Cmd handler to run
type Cmd = func(ctx context.Context) error

// Entry consists of a schedule and the func to execute on that schedule.
type Entry = cron.Entry

// Workers struct
type Workers struct {
	mu     *sync.Mutex
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	cron   *cron.Cron
}

// New returns a new workers runner.
func New(ctx context.Context) *Workers {
	w := &Workers{
		wg:   new(sync.WaitGroup),
		mu:   new(sync.Mutex),
		cron: cron.New(),
	}

	w.ctx, w.cancel = context.WithCancel(ctx)

	return w
}

// LockerOptions for exclusive running
type LockerOptions struct {
	// The maximum duration to lock a key for
	// Default: 5s
	LockTimeout time.Duration

	// The number of time the acquisition of a lock will be retried.
	// Default: 0 = do not retry
	RetryCount int

	// RetryDelay is the amount of time to wait between retries.
	// Default: 100ms
	RetryDelay time.Duration

	// Client for redis
	Client lock.RedisClient
}

// ScheduleOptions for running job
type ScheduleOptions struct {
	Name      string
	Schedule  interface{}
	Handler   Cmd
	Exclusive bool
	Locker    LockerOptions
}

// ErrEmptyName when name is empty
var ErrEmptyName = errors.New("worker must have name")

// ErrEmptyOptions when options is empty
var ErrEmptyOptions = errors.New("options must be present")

// ErrEmptyRedisClient when redis is empty
var ErrEmptyRedisClient = errors.New("redis client must be not nil")

// ErrEmptyDuration when spec or duration is empty
var ErrEmptyDuration = errors.New("spec or duration must be not nil")

// ErrEmptyHandler when handler is empty
var ErrEmptyHandler = errors.New("handler must be not null")

// The Schedule describes a job's duty cycle.
type Schedule = cron.Schedule

// Schedule adds a Job to the Cron to be run on the given schedule.
func (w *Workers) Schedule(opts *ScheduleOptions) error {
	var (
		err     error
		every   cron.Schedule
		handler func()
	)

	if len(opts.Name) == 0 {
		return ErrEmptyName
	}

	if opts == nil {
		return ErrEmptyOptions
	}

	switch schedule := opts.Schedule.(type) {
	case Schedule:
		every = schedule
	case time.Duration:
		if schedule == 0 {
			return ErrEmptyDuration
		}
		every = DelaySchedule(schedule)
	case string:
		every, err = cron.Parse(schedule)
		if err != nil {
			return wrap.Wrap(err, "cron spec parse")
		}
	default:
		return ErrEmptyDuration
	}

	if opts.Exclusive && opts.Locker.Client == nil {
		return ErrEmptyRedisClient
	}

	if opts.Handler == nil {
		return ErrEmptyHandler
	}

	handler = func() {
		if err := opts.Handler(w.ctx); err != nil {
			log.Error(wrap.Wrapf(err, "worker `%s`", opts.Name))
		}
	}

	w.cron.Schedule(every, cron.FuncJob(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("(%s) workers panic: %v", opts.Name, r)
			}
		}()

		w.add()
		defer w.done()

		if opts.Exclusive {
			if err := lock.Run(opts.Locker.Client, opts.Name, &lock.Options{
				LockTimeout: opts.Locker.LockTimeout,
				RetryCount:  opts.Locker.RetryCount,
				RetryDelay:  opts.Locker.RetryDelay,
			}, handler); err != nil {
				log.Error(wrap.Wrap(err, "locker error"))
			}
		} else {
			handler()
		}

	}))

	return nil
}

// Start the cron scheduler in its own go-routine.
func (w *Workers) Start() { w.cron.Start() }

// wait for workers stops
func (w *Workers) wait() {
	w.mu.Lock()
	w.wg.Wait()
	w.mu.Unlock()
}

// done for one of worker
func (w *Workers) done() {
	w.wg.Done()
}

// add workers to wait
func (w *Workers) add() {
	w.mu.Lock()
	w.wg.Add(1)
	w.mu.Unlock()
}

// Stop the cron scheduler and wait for jobs.
func (w *Workers) Stop() {
	log.Info("stop")
	w.cancel()
	w.cron.Stop()
	log.Info("wait")
	w.wait()
	log.Info("exit")
}

// DelaySchedule represents a simple recurring duty cycle, e.g. "Every 5 minutes".
// It does not support jobs more frequent than once a second.
type DelaySchedule time.Duration

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the second.
func (s DelaySchedule) Next(t time.Time) time.Time {
	return t.Add(time.Duration(s) - time.Duration(t.Nanosecond())/time.Millisecond)
}
