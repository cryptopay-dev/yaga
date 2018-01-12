package testdb

import (
	"fmt"
	"os"

	"github.com/go-pg/pg"
	"github.com/joho/godotenv"
)

var database *Database

type Database struct {
	DB *pg.DB
}

func GetTestDB() *Database {
	if database == nil {
		err := godotenv.Load()
		if err != nil {
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

func (d *Database) Cleanup() {
	_, err := d.DB.Exec("TRUNCATE trades, orders, users RESTART IDENTITY;")
	if err != nil {
		panic(err)
	}
}
