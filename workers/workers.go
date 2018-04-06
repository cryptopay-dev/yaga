package workers

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redis-lock"
	wrap "github.com/pkg/errors"
	"go.uber.org/atomic"
)

type TypeJob int

const (
	DefaultJob TypeJob = iota
	OnePerInstance
	OnePerCluster
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
	}

	Job func(context.Context) error

	// Options structure for creation new worker.
	Options struct {
		Name     string
		Schedule interface{}
		Handler  Job
		TypeJob  TypeJob
		Locker   LockerOptions
	}

	// Client for redis
	LockerClient = lock.RedisClient
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

	// ErrEmptyRedisClient when locker is empty
	ErrEmptyLockerClient = errors.New("locker client must be not nil")

	// ErrUnknownJobType when invalid job type
	ErrUnknownJobType = errors.New("unknown job type")
)

func (c *Cron) checkOptions(opts *Options) (Schedule, error) {
	if opts == nil {
		return nil, ErrEmptyOptions
	}
	if len(opts.Name) == 0 {
		return nil, ErrEmptyName
	}
	if opts.Handler == nil {
		return nil, ErrEmptyHandler
	}
	if opts.TypeJob == OnePerCluster && c.locker == nil {
		return nil, ErrEmptyLockerClient
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

func (c *Cron) wrapJobDefault(opts *Options) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := opts.Handler(ctx); err != nil {
			c.logger.Error(wrap.Wrapf(err, "workers `%s`", opts.Name))
		}
	}
}

func (c *Cron) wrapJobPerInstance(opts *Options) func(ctx context.Context) {
	job := c.wrapJobDefault(opts)
	lock := atomic.NewInt32(0)
	return func(ctx context.Context) {
		if !lock.CAS(0, 1) {
			return
		}
		defer lock.Store(0)
		job(ctx)
	}
}

func (c *Cron) wrapJobPerCluster(opts *Options) func(ctx context.Context) {
	job := c.wrapJobDefault(opts)
	return func(ctx context.Context) {
		if err := lock.Run(c.locker, opts.Name, &lock.Options{
			LockTimeout: opts.Locker.LockTimeout,
			RetryCount:  opts.Locker.RetryCount,
			RetryDelay:  opts.Locker.RetryDelay,
		}, func() { job(ctx) }); err != nil {
			c.logger.Error(wrap.Wrapf(err, "workers `%s` locker error", opts.Name))
		}
	}
}

// Schedule adds a Job to the Cron to be run on the given schedule.
func (c *Cron) Schedule(opts Options) error {
	schedule, err := c.checkOptions(&opts)
	if err != nil {
		return err
	}

	var job func(ctx context.Context)
	switch opts.TypeJob {
	case DefaultJob:
		job = c.wrapJobDefault(&opts)
	case OnePerInstance:
		job = c.wrapJobPerInstance(&opts)
	case OnePerCluster:
		job = c.wrapJobPerCluster(&opts)
	default:
		return ErrUnknownJobType
	}

	c.schedule(&Entry{
		Schedule: schedule,
		Name:     opts.Name,
		Job:      job,
	})

	return nil
}
