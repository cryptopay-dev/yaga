package locker

import (
	"time"

	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-redis/redis"
)

type Locker interface {
	Run(key string, timeout time.Duration, handler func() error)
}

type Lock struct {
	redis  *redis.Client
	locker *lock.Options
	logger logger.Logger
}

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

func (l *Lock) Run(key string, timeout time.Duration, handler func() error) {
	l.locker.LockTimeout = timeout
	if err := lock.Run(l.redis, key, l.locker, handler); err != nil {
		l.logger.Errorf("Locker error: %v", err.Error())
	}
}
