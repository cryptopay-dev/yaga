package testdb

import (
	"fmt"
	"strings"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/helpers/postgres"
	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
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
func GetTestDB() *Database {
	config.SetEnvPrefix("test")
	config.ReadConfig(defaultConfig)

	if database == nil {
		err := godotenv.Load()
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(err)
		}
		database = new(Database)
		database.DB, err = postgres.Connect("database")
		if err != nil {
			panic(err)
		}
	}

	return database
}
