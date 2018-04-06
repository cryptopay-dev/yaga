package redis

import (
	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/locker"
	"github.com/labstack/gommon/log"
)

// Lock struct to abstract bsm/redis-lock
type Lock struct {
	redis Client
}

// New redis locker
func New(redis Client) locker.Locker {
	return &Lock{
		redis: redis,
	}
}

// Run runs a callback handler with a Redis lock.
func (l *Lock) Run(key string, handler func(), options ...locker.Option) {
	opts := new(Options)
	opts.Parse(options...)

	if err := lock.Run(l.redis, key, &lock.Options{
		LockTimeout: opts.Timeout,
		RetryCount:  opts.RetryCount,
		RetryDelay:  opts.RetryDelay,
	}, handler); err != nil {
		log.Errorf("Locker error: %v", err.Error())
	}
}
