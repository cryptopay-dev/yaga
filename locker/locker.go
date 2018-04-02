package locker

import (
	"time"

	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/go-redis/redis"
)

// Locker interface to abstract bsm/redis-lock
type Locker interface {
	Run(key string, timeout time.Duration, handler func())
}

// Lock struct to abstract bsm/redis-lock
type Lock struct {
	redis  *redis.Client
	locker *lock.Options
}

// New creates instance of Locker
func New(opts ...Option) Locker {
	var options = newOptions(opts...)
	return &Lock{
		redis: options.Redis,
		locker: &lock.Options{
			RetryCount: 10,
			RetryDelay: 100 * time.Millisecond,
		},
	}
}

// Run runs a callback handler with a Redis lock.
func (l *Lock) Run(key string, timeout time.Duration, handler func()) {
	opts := &lock.Options{
		RetryCount:  l.locker.RetryCount,
		RetryDelay:  l.locker.RetryDelay,
		LockTimeout: timeout,
	}
	if err := lock.Run(l.redis, key, opts, handler); err != nil {
		log.Errorf("Locker error: %v", err.Error())
	}
}
