package testdb

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
)

var database *Database

// Database for tests
type Database struct {
	DB *pg.DB
}

// GetTestDB creates connection to PostgreSQL.
// Options used from ENV:
// - TEST_DATABASE_ADDR
// - TEST_DATABASE_USER
// - TEST_DATABASE_DATABASE
// - TEST_DATABASE_PASSWORD
func GetTestDB() *Database {
	if database == nil {
		err := godotenv.Load()
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(err)
		}
		database = new(Database)
		database.DB = pg.Connect(&pg.Options{
			Addr:     os.Getenv("TEST_DATABASE_ADDR"),
			User:     os.Getenv("TEST_DATABASE_USER"),
			Database: os.Getenv("TEST_DATABASE_DATABASE"),
			Password: os.Getenv("TEST_DATABASE_PASSWORD"),
			PoolSize: 2,
		})
	}

	return database
}
