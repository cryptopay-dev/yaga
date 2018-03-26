package redis

import (
	"github.com/cryptopay-dev/yaga/config"
	"github.com/go-redis/redis"
)

// Options for connections
func Options(key string) *redis.Options {
	return &redis.Options{
		Addr:        config.GetString(key + ".addr"),
		Password:    config.GetString(key + ".password"),
		DB:          config.GetInt(key + ".db"),
		PoolSize:    config.GetInt(key + ".pool_size"),
		PoolTimeout: config.GetDuration(key + ".pool_timeout"),
	}
}

// Connect to Redis and check connection:
func Connect(key string) (*redis.Client, error) {
	con := redis.NewClient(Options(key))

	// Check redis connection:
	if _, err := con.Ping().Result(); err != nil {
		return nil, err
	}

	return con, nil
}
