package redis

import (
	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/locker"
	"github.com/labstack/gommon/log"
)

// Lock struct to abstract bsm/redis-lock
type Lock struct{}

// New redis locker
func New() locker.Locker {
	return new(Lock)
}

// Run runs a callback handler with a Redis lock.
func (l *Lock) Run(key string, handler func(), options ...locker.Option) {
	opts := new(Options)
	opts.Parse(options...)

	if err := lock.Run(opts.Redis, key, &lock.Options{
		LockTimeout: opts.Timeout,
		RetryCount:  opts.RetryCount,
		RetryDelay:  opts.RetryDelay,
	}, handler); err != nil {
		log.Errorf("Locker error: %v", err.Error())
	}
}
