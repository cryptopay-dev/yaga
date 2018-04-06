package testdb

import (
	"strings"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/helpers/postgres"
	"github.com/go-pg/pg"
)

var database *Database

// Database for tests
type Database struct {
	DB *pg.DB
}

var defaultConfig = strings.NewReader(`
database:
  addr: some-address
  database: some-database
  user: some-user
  password: some-password
`)

// GetTestDB creates connection to PostgreSQL.
// Options used from ENV:
// - TEST_DATABASE_ADDR
// - TEST_DATABASE_USER
// - TEST_DATABASE_DATABASE
// - TEST_DATABASE_PASSWORD
func GetTestDB() (db *postgres.DB, err error) {
	config.SetEnvPrefix("test")

	if err = config.ReadConfig(defaultConfig); err != nil {
		return
	}

	if database == nil {
		db, err = postgres.Connect("database")
		if err != nil {
			return
		}
	}

	return db, nil
}
