package commands

import (
	"fmt"
	"reflect"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/go-pg/pg"
	"github.com/urfave/cli"
)

const configPath = "./config.yml"

type settings struct {
	Database config.Database `yaml:"database" validate:"required,dive"`
}

// FetchDB from dsn or config
func FetchDB(ctx *cli.Context, db *config.Database) (d *config.Database, err error) {
	if db != nil {
		return
	}

	dsn := ctx.String("dsn")
	if len(dsn) != 0 {
		d, err = ParseDSN(dsn)
		if err == nil {
			return
		}
	}

	var conf settings

	if err = config.Load(configPath, &conf); err == nil {
		d = &conf.Database
		return
	}

	err = fmt.Errorf("DB not set")

	return
}

// ParseDSN string to Database options
func ParseDSN(dsn string) (*config.Database, error) {
	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	return &config.Database{
		Address:  opts.Addr,
		Database: opts.Database,
		User:     opts.User,
		Password: opts.Password,
	}, err
}

// ParseConfig to Database options
func ParseConfig(i interface{}) (*config.Database, error) {
	var (
		db *config.Database
		v  = reflect.ValueOf(i)
	)

	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		if val, ok := v.Field(i).Interface().(config.Database); ok {
			db = &val
			break
		}
	}

	if db == nil {
		return nil, fmt.Errorf("config hasn't database param")
	}

	return db, nil
}
