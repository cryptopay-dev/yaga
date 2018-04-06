package redis

import (
	"time"

	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/locker"
)

// Options for creating Locker instance
type Options struct {
	Redis      Client
	RetryCount int
	RetryDelay time.Duration
	Timeout    time.Duration
}

// Parse options to locker-options
func (o *Options) Parse(opts ...locker.Option) {
	for _, op := range opts {
		op(o)
	}
}

// Client is a minimal client interface.
type Client = lock.RedisClient

// Timeout closure
func Timeout(v time.Duration) locker.Option {
	return func(o locker.Options) {
		o.(*Options).Timeout = v
	}
}

// Count closure to set RetryCount
func Count(v int) locker.Option {
	return func(o locker.Options) {
		o.(*Options).RetryCount = v
	}
}

// Delay closure to set RetryDelay
func Delay(v time.Duration) locker.Option {
	return func(o locker.Options) {
		o.(*Options).RetryDelay = v
	}
}
