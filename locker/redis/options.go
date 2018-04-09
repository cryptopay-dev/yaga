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
func (o *Options) Parse(opts ...locker.Option) error {
	for _, op := range opts {
		if err := op(o); err != nil {
			return err
		}
	}
	return nil
}

// Client is a minimal client interface.
type Client = lock.RedisClient

// ErrUnknownOptions passed to Option
var ErrUnknownOptions = errors.New("unknown options passed")

// Timeout closure
func Timeout(v time.Duration) locker.Option {
	return func(o locker.Options) error {
		if op, ok := o.(*Options); ok {
			op.Timeout = v
			return nil
		}
		return ErrUnknownOptions
	}
}

// Count closure to set RetryCount
func Count(v int) locker.Option {
	return func(o locker.Options) error {
		if op, ok := o.(*Options); ok {
			op.RetryCount = v
			return nil
		}
		return ErrUnknownOptions
	}
}

// Delay closure to set RetryDelay
func Delay(v time.Duration) locker.Option {
	return func(o locker.Options) error {
		if op, ok := o.(*Options); ok {
			op.RetryDelay = v
			return nil
		}
		return ErrUnknownOptions
	}
}
