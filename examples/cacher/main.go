package main

import (
	"github.com/cryptopay-dev/yaga/cacher"
	"github.com/cryptopay-dev/yaga/cacher/redis"
	redisStore "github.com/go-redis/redis"
)

func variantOne() cacher.Cacher {
	var r = redisStore.NewClient(&redisStore.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	return redis.FromConnection(r)
}

func variantTwo() cacher.Cacher {
	return redis.New(
		redis.Address("127.0.0.1:6379"),
		redis.Password(""),
		redis.DB(0),
	)
}
