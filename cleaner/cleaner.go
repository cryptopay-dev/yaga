package cleaner

import (
	"fmt"
	"strings"
	"time"

	"github.com/cryptopay-dev/yaga/cacher"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-redis/redis"
)

const (
	infStartOrderbookCleaner      = "Starting orderbooks-cleaner"
	infTryToUpdateOrderbookTTLTpl = "Try to update orderbook TTL for '%s':'%s' from %s to %s"

	errFailUpdateOrderbookTTLTpl   = "fail update orderbook TTL for ['%s'/'%s']: %v"
	errCannotFetchOrderbookKeysTpl = "can't fetch orderbooks-keys from redis: %v"

	LockerItemTpl = "midas-fx:%s:orderbook:%s"
)

type Pair struct {
	Name     string
	Expire   time.Duration
	Platform string
}

type Pairs = []Pair

type RedisData = map[string]time.Duration

type service struct {
	Pairs  Pairs
	Cacher cacher.Cacher
	Logger logger.Logger
}

func preparePairs(conf Platforms) Pairs {
	var result = make(Pairs, 0)

	for platform, platformConfig := range conf {
		for _, pair := range platformConfig.Pairs {
			result = append(result, Pair{
				Platform: platform,
				Name:     pair.InternalName,
				Expire:   platformConfig.ExpiresEvery,
			})
		}
	}

	return result
}

func New(opts ...Option) Cleaner {
	var options = newOptions(opts...)
	return &service{
		Cacher: options.Cacher,
		Logger: options.Logger,
		Pairs:  preparePairs(options.Platforms),
	}
}

func (c *service) prepareRedisData(pattern string) (RedisData, error) {
	var (
		err    error
		result []string
		data   = make(RedisData)
	)
	if result, err = c.Cacher.Keys(pattern); err != nil && err != redis.Nil {
		return nil, err
	}

	for _, key := range result {
		var ttl time.Duration

		if ttl, err = c.Cacher.TTL(key); err != nil && err != redis.Nil {
			return nil, err
		}

		data[key] = ttl
	}

	return data, nil
}

func fmtOrderbookKey(platform, pair string) string {
	return fmt.Sprintf(
		LockerItemTpl,
		strings.ToUpper(platform),
		strings.ToUpper(pair),
	)
}

func (c *service) UpdateTTL() error {
	var (
		err     error
		data    RedisData
		pattern = fmtOrderbookKey("*", "*")
	)

	c.Logger.Info(infStartOrderbookCleaner)
	if data, err = c.prepareRedisData(pattern); err != nil {
		return fmt.Errorf(errCannotFetchOrderbookKeysTpl, err)
	}

	for _, pair := range c.Pairs {
		var (
			ok  bool
			ttl time.Duration
			key = fmtOrderbookKey(pair.Platform, pair.Name)
		)

		if ttl, ok = data[key]; !ok {
			continue
		}

		if ttl > pair.Expire {
			c.Logger.Infof(infTryToUpdateOrderbookTTLTpl, pair.Platform, pair.Name, ttl, pair.Expire)
			if err = c.Cacher.Expire(key, pair.Expire); err != nil {
				return fmt.Errorf(errFailUpdateOrderbookTTLTpl, pair.Platform, pair.Name, err)
			}
		}
	}

	return nil
}
