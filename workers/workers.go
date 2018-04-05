package workers

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redis-lock"
)

type (
	// LockerOptions for exclusive running
	LockerOptions struct {
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

	// Options structure for creation new worker.
	Options struct {
		Name      string
		Schedule  interface{}
		Handler   func(context.Context) error
		Exclusive bool
		Locker    LockerOptions
	}
)

var (
	// ErrEmptyOptions when options is empty
	ErrEmptyOptions = errors.New("options must be present")

	// ErrEmptyName when name is empty
	ErrEmptyName = errors.New("worker must have name")

	// ErrEmptyHandler when handler is empty
	ErrEmptyHandler = errors.New("handler must be not null")

	// ErrEmptyDuration when spec or duration is empty
	ErrEmptyDuration = errors.New("spec or duration must be not nil")

	// ErrEmptyRedisClient when redis is empty
	ErrEmptyRedisClient = errors.New("redis client must be not nil")
)

type Workers struct {
	cron   *Cron
	ctx    context.Context
	cancel context.CancelFunc
}

func New(ctx context.Context) *Workers {
	ctx, cancel := context.WithCancel(ctx)
	return &Workers{
		ctx:    ctx,
		cancel: cancel,
		cron:   NewCron(),
	}
}

func (w *Workers) Start() {
	w.cron.Start(w.ctx)
}

func (w *Workers) Stop() {
	w.cancel()
}

func (w *Workers) Wait(ctx context.Context) error {
	return w.cron.Wait(ctx)
}

func checkOptions(opts *Options) (Schedule, error) {
	if opts == nil {
		return nil, ErrEmptyOptions
	}
	if len(opts.Name) == 0 {
		return nil, ErrEmptyName
	}
	if opts.Handler == nil {
		return nil, ErrEmptyHandler
	}
	if opts.Exclusive && opts.Locker.Client == nil {
		return nil, ErrEmptyRedisClient
	}

	var err error
	var schedule Schedule
	switch sc := opts.Schedule.(type) {
	case string:
		schedule, err = Parse(sc)
		if err != nil {
			return nil, err
		}
	case time.Duration:
		if sc == 0 {
			return nil, ErrEmptyDuration
		}
		schedule = Every(sc)
	case ConstantDelaySchedule:
		schedule = sc
	case DelaySchedule:
		schedule = sc
	default:
		return nil, ErrEmptyDuration
	}

	return schedule, nil
}

func (w *Workers) Schedule(opts Options) error {
	schedule, err := checkOptions(&opts)
	if err != nil {
		return err
	}

	/*
		if opts.Exclusive {
			job = func(ctx context.Context) {
				if err := lock.Run(opts.Locker.Client, opts.Name, &lock.Options{
					LockTimeout: opts.Locker.LockTimeout,
					RetryCount:  opts.Locker.RetryCount,
					RetryDelay:  opts.Locker.RetryDelay,
				}, func() { handler(ctx) }); err != nil {
					log.Error(wrap.Wrap(err, "locker error"))
				}
			}
		}
	*/

	w.cron.Schedule(schedule, opts.Name, opts.Handler)

	return nil
}

// DelaySchedule represents a simple recurring duty cycle, e.g. "Every 5 minutes".
// It does not support jobs more frequent than once a millisecond.
type DelaySchedule time.Duration

// Next returns the next time this should be run.
// This rounds so that the next activation time will be on the millisecond.
func (s DelaySchedule) Next(t time.Time) time.Time {
	d := time.Duration(s) - time.Duration(t.Nanosecond())/time.Millisecond
	if d < time.Millisecond {
		d = time.Millisecond
	}
	return t.Add(d)
}
