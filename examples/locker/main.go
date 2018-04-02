package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cryptopay-dev/yaga/locker"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/go-redis/redis"
)

func main() {
	store := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	log.Init()

	lock := locker.New(
		locker.Redis(store),
	)

	wg := sync.WaitGroup{}

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func(index int) {
			lock.Run("my-key", time.Second*10, func() {
				// Do some thing ... //

				// For example
				fmt.Println("Step :", index)
				wg.Done()
			})
		}(i)
	}

	wg.Wait()
}
