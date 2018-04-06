package workers

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/atomic"
)

var (
	// ErrAlreadyWorker is returned by New calls
	// when workers name is already exists.
	ErrAlreadyWorker = errors.New("worker name must be unique")

	_ Locker = &LockerJobPerInstance{}
	_ Locker = &LockerUniqJobPerInstance{}
)

type LockerJobPerInstance struct{}

func (LockerJobPerInstance) TypeJob() TypeJob {
	return OnePerInstance
}

func (LockerJobPerInstance) WrapJob(opts interface{}, jobName string, job LockerJob) (LockerJob, error) {
	lock := atomic.NewInt32(0)
	return func(ctx context.Context) {
		if !lock.CAS(0, 1) {
			return
		}
		defer lock.Store(0)
		job(ctx)
	}, nil
}

type LockerUniqJobPerInstance struct {
	LockerJobPerInstance

	mu    sync.Mutex
	names map[string]struct{}
}

func (l *LockerUniqJobPerInstance) WrapJob(opts interface{}, jobName string, job LockerJob) (LockerJob, error) {
	l.mu.Lock()
	if l.names == nil {
		l.names = make(map[string]struct{})
	}
	if _, found := l.names[jobName]; found {
		l.mu.Unlock()
		return nil, ErrAlreadyWorker
	}
	l.mu.Unlock()

	return l.LockerJobPerInstance.WrapJob(opts, jobName, job)
}
