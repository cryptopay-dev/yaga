package redis

import (
	"time"

	"github.com/cryptopay-dev/yaga/cacher"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type service struct {
	redis *redis.Client
}

// New creates redis-cache instance
func New(opts ...Option) cacher.Cacher {
	var options = newOptions(opts...)

	return &service{
		redis: redis.NewClient(&redis.Options{
			Addr:     options.Address,
			Password: options.Password,
			DB:       options.DB,
		}),
	}
}

// FromConnection creates redis-cache from existing redis-connection
func FromConnection(client *redis.Client) cacher.Cacher {
	return &service{redis: client}
}

// Get key from redis-cache to val interface
func (s *service) Get(key string, val interface{}) error {
	if err := s.redis.Get(key).Scan(val); err != nil && err != redis.Nil {
		return errors.Wrap(err, "cacher Get failed")
	}
	return nil
}

// Set key to redis-cache with value, and expiration
func (s *service) Set(key string, val interface{}, expiration time.Duration) error {
	return errors.Wrap(s.redis.Set(key, val, expiration).Err(), "cacher Set failed")
}

// Del key from redis-cache
func (s *service) Del(key string) error {
	return errors.Wrap(s.redis.Del(key).Err(), "cacher Del failed")
}

// Exists check keys in redis-cache
func (s *service) Exists(keys ...string) (n int64, err error) {
	n, err = s.redis.Exists(keys...).Result()
	return n, errors.Wrap(err, "cacher Exists failed")
}

// Keys fetch from redis-cache by pattern
func (s *service) Keys(pattern string) (keys []string, err error) {
	keys, err = s.redis.Keys(pattern).Result()
	return keys, errors.Wrap(err, "cacher Keys failed")
}

// TTL fetch for key from redis-cache
func (s *service) TTL(key string) (d time.Duration, err error) {
	d, err = s.redis.TTL(key).Result()
	return d, errors.Wrap(err, "cacher TTL failed")
}

// Expire sets for key in redis-cache
func (s *service) Expire(key string, duration time.Duration) error {
	return errors.Wrap(s.redis.Expire(key, duration).Err(), "cacher Expire failed")
}
