package cli

import (
	"context"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
	"github.com/go-redis/redis"
)

// RunOptions for pass db, redis, etc to application:
type RunOptions struct {
	DB           *pg.DB
	Redis        *redis.Client
	Logger       logger.Logger
	Debug        bool
	BuildTime    string
	BuildVersion string
}

// Instance abstraction layer above Application
type Instance interface {
	Run(RunOptions) error
	Shutdown(ctx context.Context) error
}
