package internal

import (
	"errors"
	"reflect"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/go-pg/pg"
)

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
		return nil, errors.New("config hasn't database param")
	}

	return db, nil
}
