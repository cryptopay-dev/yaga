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
	"go.uber.org/atomic"
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
	cron   *cron.Cron
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
	jobCh  chan func()
	state  *atomic.Int32
}

func New(ctx context.Context) *Workers {
	w := &Workers{
		cron:  cron.New(),
		done:  make(chan struct{}),
		jobCh: make(chan func()),
		state: atomic.NewInt32(0),
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	go w.dispatcher()

	return w
}

func (w *Workers) Start() {
	if w.state.CAS(0, 1) {
		w.jobCh <- w.cron.Run
	}
}

func (w *Workers) Stop() {
	if w.state.CAS(1, 2) {
		w.cancel()
	}
}

func (w *Workers) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.done:
	}

	return nil
}

func (w *Workers) checkOptions(opts *Options) (cron.Schedule, error) {
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
	var schedule cron.Schedule
	switch sc := opts.Schedule.(type) {
	case string:
		schedule, err = cron.Parse(sc)
		if err != nil {
			return nil, err
		}
	case time.Duration:
		if sc == 0 {
			return nil, ErrEmptyDuration
		}
		schedule = cron.Every(sc)
	case DelaySchedule:
		schedule = sc
	default:
		return nil, ErrEmptyDuration
	}

	return schedule, nil
}

func (w *Workers) recovery(workName string) {
	if r := recover(); r != nil {
		log.Errorf("workers `%s` panic: %v", workName, r)
	}
}

func (w *Workers) Schedule(opts Options) error {
	schedule, err := w.checkOptions(&opts)
	if err != nil {
		return err
	}
	handler := func() {
		defer w.recovery(opts.Name)
		if err := opts.Handler(w.ctx); err != nil {
			log.Error(wrap.Wrapf(err, "worker `%s`", opts.Name))
		}
	}
	job := handler

	if opts.Exclusive {
		job = func() {
			if err := lock.Run(opts.Locker.Client, opts.Name, &lock.Options{
				LockTimeout: opts.Locker.LockTimeout,
				RetryCount:  opts.Locker.RetryCount,
				RetryDelay:  opts.Locker.RetryDelay,
			}, handler); err != nil {
				log.Error(wrap.Wrap(err, "locker error"))
			}
		}
	}

	w.cron.Schedule(schedule, cron.FuncJob(func() {
		select {
		case w.jobCh <- job:
		case <-w.ctx.Done():
		}
	}))

	return nil
}

func (w *Workers) dispatcher() {
	wg := new(sync.WaitGroup)
	for {
		select {
		case job := <-w.jobCh:
			wg.Add(1)
			go func() {
				defer wg.Done()
				job()
			}()
		case <-w.ctx.Done():
			w.cron.Stop()
			wg.Wait()
			close(w.done)
			return
		}
	}
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
