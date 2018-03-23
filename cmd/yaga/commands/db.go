package commands

import (
	"github.com/cryptopay-dev/yaga/config"
	"github.com/go-pg/pg"
	"github.com/urfave/cli"
)

const configPath = "./config.yml"

// FetchDB from dsn or config
func FetchDB(ctx *cli.Context, key string) (err error) {
	dsn := ctx.String("dsn")

	if len(dsn) != 0 {
		return
	}

	return ParseDSN(key, dsn)
}

// ParseDSN string to Database options
func ParseDSN(key, dsn string) error {
	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return err
	}

	config.Set(key+".addr", opts.Addr)
	config.Set(key+".user", opts.User)
	config.Set(key+".password", opts.Password)
	config.Set(key+".database", opts.Database)

	return nil
}
