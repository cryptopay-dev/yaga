package cacher

import "time"

// Cacher interface to abstract Redis/other
type Cacher interface {
	Get(key string, val interface{}) error
	Set(key string, val interface{}, timeout time.Duration) error
	Keys(pattern string) ([]string, error)
	TTL(key string) (time.Duration, error)
	Expire(key string, duration time.Duration) error
	Del(key string) error
	Exists(key string) (int64, error)
}
