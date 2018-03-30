package postgres

import (
	"github.com/cryptopay-dev/yaga/config"
	"github.com/go-pg/pg"
)

// Options for connections
func Options(key string) *pg.Options {
	return &pg.Options{
		Addr:     config.GetString(key + ".addr"),
		User:     config.GetString(key + ".user"),
		Password: config.GetString(key + ".password"),
		Database: config.GetString(key + ".database"),
	}
}

// DB for hide pg.DB
type DB = pg.DB

// Connect to PostgreSQL and check connection:
func Connect(key string) (*DB, error) {
	con := pg.Connect(Options(key))

	// Check postgres connection:
	if _, err := con.ExecOne("SELECT 1"); err != nil {
		return nil, err
	}

	return con, nil
}
