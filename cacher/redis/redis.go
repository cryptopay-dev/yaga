package redis

import (
	"time"

	"github.com/cryptopay-dev/yaga/cacher"
	"github.com/go-redis/redis"
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
		return err
	}
	return nil
}

// Set key to redis-cache with value, and expiration
func (s *service) Set(key string, val interface{}, expiration time.Duration) error {
	return s.redis.Set(key, val, expiration).Err()
}

// Del key from redis-cache
func (s *service) Del(key string) error {
	return s.redis.Del(key).Err()
}

// Exists check keys in redis-cache
func (s *service) Exists(keys ...string) (int64, error) {
	return s.redis.Exists(keys...).Result()
}

// Keys fetch from redis-cache by pattern
func (s *service) Keys(pattern string) ([]string, error) {
	return s.redis.Keys(pattern).Result()
}

// TTL fetch for key from redis-cache
func (s *service) TTL(key string) (time.Duration, error) {
	return s.redis.TTL(key).Result()
}

// Expire sets for key in redis-cache
func (s *service) Expire(key string, duration time.Duration) error {
	return s.redis.Expire(key, duration).Err()
}
