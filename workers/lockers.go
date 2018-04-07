// TODO move to packet yaga/locker(/...)
package workers

import (
	"context"
	"errors"

	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/logger/log"
	"go.uber.org/atomic"
)

// TODO need to discuss about context.Context
type LockerHandler func(context.Context)

type Locker interface {
	// TODO move to packet yaga/locker
	// while Options as interface{} (may be later as type Options)
	// and while as Wrapper
	WrapRun(key string, options interface{}, handler LockerHandler) (LockerHandler, error)
}

var (
	_ Locker = &LockerJobPerInstance{}

	// ErrUnknownOptions passed to Option
	ErrUnknownOptions = errors.New("unknown options passed")
)

type LockerJobPerInstance struct{}

func (LockerJobPerInstance) WrapRun(key string, opts interface{}, handler LockerHandler) (LockerHandler, error) {
	// TODO need verification unique for key?
	lock := atomic.NewInt32(0)
	return func(ctx context.Context) {
		if !lock.CAS(0, 1) {
			return
		}
		defer lock.Store(0)
		handler(ctx)
	}, nil
}

// *******************
// TODO example locker for OnePerCluster
// *******************

// LockerJobPerCluster struct for OnePerCluster job type
type LockerJobPerCluster struct {
	redis Client
}

// Client is a minimal client interface.
type Client = lock.RedisClient

// New locker
// TODO while 'NewLocker'
func NewLocker(client Client) Locker {
	return &LockerJobPerCluster{
		redis: client,
	}
}

// TODO while 'LockerOptions'
type LockerOptions = lock.Options

func (l *LockerJobPerCluster) WrapRun(key string, opts interface{}, handler LockerHandler) (LockerHandler, error) {
	options, ok := opts.(LockerOptions)
	if !ok {
		return nil, ErrUnknownOptions
	}

	return func(ctx context.Context) {
		if err := lock.Run(l.redis, key, &options, func() {
			handler(ctx)
		}); err != nil {
			log.Debugf("Locker error: %v", err)
		}
	}, nil
}
