package config

import (
	"github.com/go-redis/redis"
)

// Redis default configuration
type Redis struct {
	Address  string `yaml:"address" validate:"required"`
	Database int    `yaml:"database"`
	Password string `yaml:"password"`
}

// Connect to Redis and check connection:
func (r Redis) Connect() (*redis.Client, error) {
	cache := redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.Database,
	})

	// Check redis connection:
	if _, err := cache.Ping().Result(); err != nil {
		return nil, err
	}

	return cache, nil
}
