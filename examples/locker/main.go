package main

import (
	"fmt"
	"sync"
	"time"

	locker "github.com/cryptopay-dev/yaga/locker/redis"
	"github.com/go-redis/redis"
)

func main() {
	store := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	lock := locker.New(store)

	wg := sync.WaitGroup{}

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func(index int) {
			lock.Run("my-key", func() {
				// Do some thing ... //

				// For example
				fmt.Println("Step :", index)
				wg.Done()
			}, locker.Timeout(time.Second*10))
		}(i)
	}

	wg.Wait()
}
