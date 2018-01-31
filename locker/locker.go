package locker

import (
	"time"

	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-redis/redis"
)

// Locker interface to abstract bsm/redis-lock
type Locker interface {
	Run(key string, timeout time.Duration, handler func() error)
}

// Lock struct to abstract bsm/redis-lock
type Lock struct {
	redis  *redis.Client
	locker *lock.Options
	logger logger.Logger
}

// New creates instance of Locker
func New(opts ...Option) Locker {
	var options = newOptions(opts...)
	return &Lock{
		redis:  options.Redis,
		logger: options.Logger,
		locker: &lock.Options{
			RetryCount: 10,
			RetryDelay: 100 * time.Millisecond,
		},
	}
}

// Run runs a callback handler with a Redis lock.
func (l *Lock) Run(key string, timeout time.Duration, handler func() error) {
	l.locker.LockTimeout = timeout
	if err := lock.Run(l.redis, key, l.locker, handler); err != nil {
		l.logger.Errorf("Locker error: %v", err.Error())
	}
}
