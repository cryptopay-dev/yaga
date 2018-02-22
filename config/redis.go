package config

import (
	"time"

	"github.com/go-redis/redis"
)

// Redis default configuration
type Redis struct {
	Address     string        `yaml:"address" validate:"required"`
	Database    int           `yaml:"database"`
	Password    string        `yaml:"password"`
	PoolSize    int           `yaml:"pool_size" validate:"required,gte=0"`
	PoolTimeout time.Duration `yaml:"pool_timeout" validate:"required,gte=0"`
}

// Connect to Redis and check connection:
func (r Redis) Connect() (*redis.Client, error) {
	cache := redis.NewClient(&redis.Options{
		Addr:        r.Address,
		Password:    r.Password,
		DB:          r.Database,
		PoolSize:    r.PoolSize,
		PoolTimeout: r.PoolTimeout,
	})

	// Check redis connection:
	if _, err := cache.Ping().Result(); err != nil {
		return nil, err
	}

	return cache, nil
}
