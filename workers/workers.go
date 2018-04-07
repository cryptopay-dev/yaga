package workers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/cryptopay-dev/yaga/logger/log"
	wrap "github.com/pkg/errors"
	"go.uber.org/atomic"
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
)

var (
	// ErrAlreadyWorker is returned by New calls
	// when workers name is already exists.
	ErrAlreadyWorker = errors.New("worker name must be unique")

	// ErrEmptyOptions when options is empty
	ErrEmptyOptions = errors.New("options must be present")

	// ErrEmptyName when name is empty
	ErrEmptyName = errors.New("worker must have name")

	// ErrEmptyHandler when handler is empty
	ErrEmptyHandler = errors.New("handler must be not null")

	// ErrEmptyDuration when spec or duration is empty
	ErrEmptyDuration = errors.New("spec or duration must be not nil")

	// ErrUnknownJobType when invalid job type
	ErrUnknownJobType = errors.New("unknown job type")
)

// Workers keeps track of any number of entries, invoking the associated func as
// specified by the schedule. It may be started, stopped, and the entries may
// be inspected while running.
type Workers struct {
	entries []*entry
	done    chan struct{}
	add     chan *entry
	state   *atomic.Int32
	logger  logger.Logger
	lockers map[TypeJob]Locker

	mu    *sync.Mutex
	names map[string]struct{}
}

// New returns a new Workers job runner, in the Local time zone.
func New(lockerOnePerCluster Locker) *Workers {
	w := &Workers{
		entries: nil,
		add:     make(chan *entry),
		done:    make(chan struct{}),
		state:   atomic.NewInt32(0),
		logger:  log.Logger(),
		lockers: make(map[TypeJob]Locker),
		mu:      new(sync.Mutex),
		names:   make(map[string]struct{}),
	}

	w.lockers[OnePerInstance] = new(LockerJobPerInstance)
	if lockerOnePerCluster != nil {
		w.lockers[OnePerCluster] = lockerOnePerCluster
	}

	return w
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

	if opts.TypeJob > maxTypeJob {
		return nil, ErrUnknownJobType
	}
	var err error
	var schedule Schedule
	switch sc := opts.Schedule.(type) {
	case string:
		schedule, err = parse(sc)
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

func (w *Workers) regJobName(name string) error {
	w.mu.Lock()
	if _, found := w.names[name]; found {
		w.mu.Unlock()
		return ErrAlreadyWorker
	}
	w.names[name] = struct{}{}
	w.mu.Unlock()

	return nil
}

// Schedule adds a Job to the Workers to be run on the given schedule.
func (w *Workers) Schedule(opts Options) error {
	schedule, err := checkOptions(&opts)
	if err != nil {
		return err
	}

	job := func(ctx context.Context) {
		if err := opts.Handler(ctx); err != nil {
			w.logger.Error(wrap.Wrapf(err, "workers `%s`", opts.Name))
		}
	}

	if locker, ok := w.lockers[opts.TypeJob]; ok {
		// TODO wait until Locker interface will be approved
		job, err = locker.WrapRun(opts.Name, opts.Locker, job)
		if err != nil {
			return err
		}
	}

	// TODO for all job types or only for one job type?
	if err = w.regJobName(opts.Name); err != nil {
		return err
	}

	w.schedule(&entry{
		Schedule: schedule,
		Name:     opts.Name,
		Job:      job,
	})

	return nil
}
