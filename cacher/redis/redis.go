package redis

import (
	"time"

	"github.com/cryptopay-dev/yaga/cacher"
	"github.com/go-redis/redis"
)

type service struct {
	redis *redis.Client
}

func New(opts ...Option) cacher.Cacher {
	var options = newOptions(opts...)
	_ = options

	return &service{
		redis: redis.NewClient(&redis.Options{
			Addr:     options.Address,
			Password: options.Password,
			DB:       options.DB,
		}),
	}
}

func FromConnection(client *redis.Client) cacher.Cacher {
	return &service{redis: client}
}

func (s *service) Get(key string, val interface{}) error {
	if err := s.redis.Get(key).Scan(val); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func (s *service) Set(key string, val interface{}, expiration time.Duration) error {
	return s.redis.Set(key, val, expiration).Err()
}

func (s *service) Del(key string) error {
	return s.redis.Del(key).Err()
}

func (s *service) Keys(pattern string) ([]string, error) {
	return s.redis.Keys(pattern).Result()
}

func (s *service) TTL(key string) (time.Duration, error) {
	return s.redis.TTL(key).Result()
}

func (s *service) Expire(key string, duration time.Duration) error {
	return s.redis.Expire(key, duration).Err()
}
