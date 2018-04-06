package workers

import (
	"context"
	"errors"
	"time"

	wrap "github.com/pkg/errors"
)

type TypeJob int

const (
	DefaultJob TypeJob = iota
	OnePerInstance
	OnePerCluster
)

var maxTypeJob = OnePerCluster

type (
	Job func(context.Context) error

	// Options structure for creation new worker.
	Options struct {
		Name     string
		Schedule interface{}
		Handler  Job
		TypeJob  TypeJob
		Locker   interface{}
	}

	LockerJob func(context.Context)

	Locker interface {
		TypeJob() TypeJob
		WrapJob(lockerOptions interface{}, jobName string, job LockerJob) (LockerJob, error)
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

	// ErrEmptyLocker when locker is empty
	ErrEmptyLocker = errors.New("locker must be not nil")

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

	if opts.TypeJob > maxTypeJob {
		return nil, ErrUnknownJobType
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
	case Schedule:
		schedule = sc
	default:
		return nil, ErrEmptyDuration
	}

	return schedule, nil
}

// Schedule adds a Job to the Cron to be run on the given schedule.
func (c *Cron) Schedule(opts Options) error {
	schedule, err := c.checkOptions(&opts)
	if err != nil {
		return err
	}

	job := func(ctx context.Context) {
		if err := opts.Handler(ctx); err != nil {
			c.logger.Error(wrap.Wrapf(err, "workers `%s`", opts.Name))
		}
	}

	if locker, ok := c.lockers[opts.TypeJob]; ok {
		job, err = locker.WrapJob(opts.Locker, opts.Name, job)
		if err != nil {
			return err
		}
	}

	c.schedule(&Entry{
		Schedule: schedule,
		Name:     opts.Name,
		Job:      job,
	})

	return nil
}
