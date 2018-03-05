package redis

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cryptopay-dev/yaga/cacher"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

const keyTpl = "test:key:%d"

func closedCacher() cacher.Cacher {
	c := redis.NewClient(&redis.Options{
		Addr: os.Getenv("TEST_REDIS_ADDR"),
	})
	defer c.Close()
	return FromConnection(c)
}

func defaultCacher() cacher.Cacher {
	return New(
		Address(os.Getenv("TEST_REDIS_ADDR")),
		DB(0),
		Password(""),
	)
}

func TestService(t *testing.T) {
	c := defaultCacher()

	var (
		i    = 1
		val  int
		key  = fmt.Sprintf(keyTpl, 0)
		ttl1 = time.Second * time.Duration(100)
		ttl3 = time.Second * time.Duration(105)
	)

	t.Run("cacher.Set", func(t *testing.T) {
		// Try to cache value:
		if err := c.Set(key, i, ttl1); !assert.NoError(t, err) {
			t.FailNow()
		}
	})

	t.Run("cacher.Get", func(t *testing.T) {
		// Try to get value from cache:
		if err := c.Get(key, &val); !assert.NoError(t, err) || !assert.Equal(t, i, val) {
			t.FailNow()
		}

		assert.NotPanics(t, func() {
			// Try to get value from cache (error):
			if err := closedCacher().Get("", &val); !assert.Error(t, err) {
				t.FailNow()
			}
		})
	})

	t.Run("cacher.TTL", func(t *testing.T) {
		// Try to get ttl key from cache:
		if ttl2, err := c.TTL(key); !assert.NoError(t, err) || !assert.Equal(t, ttl1, ttl2) {
			t.FailNow()
		}
	})

	t.Run("cacher.Expire", func(t *testing.T) {
		// Try to set ttl key in cache:
		if err := c.Expire(key, ttl3); !assert.NoError(t, err) {
			t.FailNow()
		}

		// Try to get ttl key from cache:
		if ttl2, err := c.TTL(key); !assert.NoError(t, err) || !assert.Equal(t, ttl3, ttl2) {
			t.FailNow()
		}
	})

	t.Run("cacher.Exists", func(t *testing.T) {
		// Try to check exists key:
		if count, err := c.Exists(key); !assert.NoError(t, err) || !assert.True(t, count == 1) {
			t.FailNow()
		}
	})

	t.Run("cacher.Keys", func(t *testing.T) {
		// Try to get keys from cache:
		if keys, err := c.Keys(key); !assert.NoError(t, err) || !assert.Equal(t, []string{key}, keys) {
			t.FailNow()
		}
	})

	t.Run("cacher.Del", func(t *testing.T) {
		// Try to remove key from cache:
		if err := c.Del(key); !assert.NoError(t, err) {
			t.FailNow()
		}

		// Try to check exists key:
		if count, err := c.Exists(key); !assert.NoError(t, err) || !assert.True(t, count == 0) {
			t.FailNow()
		}
	})
}
