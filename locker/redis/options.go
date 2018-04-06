package redis

import (
	"errors"
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
func Timeout(v interface{}) locker.Option {
	return func(o locker.Options) error {
		val, ok := v.(time.Duration)
		if !ok {
			return errors.New("bad option timeout, must be time.Duration")
		}
		o.(*Options).Timeout = val
		return nil
	}
}

// Count closure to set RetryCount
func Count(v interface{}) locker.Option {
	return func(o locker.Options) error {
		val, ok := v.(int)
		if !ok {
			return errors.New("bad option count, must be int")
		}
		o.(*Options).RetryCount = val
		return nil
	}
}

// Delay closure to set RetryDelay
func Delay(v interface{}) locker.Option {
	return func(o locker.Options) error {
		val, ok := v.(time.Duration)
		if !ok {
			return errors.New("bad option timeout, must be time.Duration")
		}
		o.(*Options).RetryDelay = val
		return nil
	}
}

// Redis closure to set field in Options
func Redis(v interface{}) locker.Option {
	return func(o locker.Options) error {
		val, ok := v.(Client)
		if !ok {
			return errors.New("bad option redis, must be RedisClient")
		}
		o.(*Options).Redis = val
		return nil
	}
}
