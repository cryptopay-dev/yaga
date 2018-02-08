package cli

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
)

// RunOptions for pass db, redis, etc to application:
type RunOptions struct {
	DB    *pg.DB
	Redis *redis.Client
}

// Instance abstraction layer above Application
type Instance interface {
	Run(RunOptions) error
	Shutdown(ctx context.Context) error
}
