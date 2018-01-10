package cleaner

import (
	"testing"
	"time"

	"github.com/cryptopay-dev/yaga/cacher"
	"github.com/cryptopay-dev/yaga/cacher/redis"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/stretchr/testify/assert"
)

const (
	fakePair     = "fake-pair"
	fakeValue    = "fakeValue"
	fakeMapping  = "fake-mapping"
	fakePlatform = "fake-platform"
)

func testConfig(ttl time.Duration) Platforms {
	return Platforms{
		fakePlatform: PlatformConfig{
			Pairs: []PairConfig{
				{
					InternalName: fakePair,
					Mapping:      fakeMapping,
				},
			},
			CollectEvery: ttl,
			ExpiresEvery: ttl,
		},
	}

}

func testCleaner(ttl time.Duration, cli cacher.Cacher) Cleaner {
	var conf = testConfig(ttl)
	return New(
		Platform(conf),
		Logger(nop.New()),
		Cacher(cli),
	)
}

func testRedis() cacher.Cacher {
	return redis.New(
		redis.DB(9),
		redis.Address("127.0.0.1:6379"),
		redis.Password(""),
	)
}

func TestCleaner(t *testing.T) {
	var err error

	for i := 100; i > 0; i -= 10 {
		var (
			ttl         time.Duration
			dur         = time.Second * time.Duration(i)
			key         = fmtOrderbookKey(fakePlatform, fakePair)
			pat         = fmtOrderbookKey("*", "*")
			cli         = testRedis()
			clean       = testCleaner(dur, cli)
			expectedTTL = time.Second * time.Duration(i)
		)

		cli.Set(key, fakeValue, time.Second*time.Duration(i+10))

		err = clean.UpdateTTL()
		assert.NoError(t, err)

		ttl, err = cli.TTL(key)
		assert.NoError(t, err)
		assert.Equal(t, expectedTTL, ttl)

		data, err := clean.(*service).prepareRedisData(pat)
		assert.NoError(t, err)
		for name, durTTL := range data {
			assert.Equal(t, key, name)
			assert.Equal(t, ttl, durTTL)
		}
	}
}
