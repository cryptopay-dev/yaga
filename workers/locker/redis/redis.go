package redis

import (
	"github.com/bsm/redis-lock"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/cryptopay-dev/yaga/workers/locker"
)

// *******************
// TODO example locker for OnePerCluster
// *******************

// Lock struct for OnePerCluster job type
type Lock struct {
	redis Client
}

// Client is a minimal client interface.
type Client = lock.RedisClient

// New locker
func New(client Client) *Lock {
	return &Lock{
		redis: client,
	}
}

type Options = lock.Options

type lockItem struct {
	redis   Client
	options *Options
}

func (l *Lock) GetLocker(opts Options) locker.Locker {
	return &lockItem{
		redis:   l.redis,
		options: &opts,
	}
}

// TODO need returns error?
func (l *lockItem) Run(key string, handler func()) {
	if err := lock.Run(l.redis, key, l.options, handler); err != nil {
		log.Debugf("Locker '%s' error: %v", key, err)
	}
}
