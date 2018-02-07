package config

import (
	"github.com/go-pg/pg"
)

// Database base configuration:
type Database struct {
	Address  string `yaml:"address" validate:"required"`
	Database string `yaml:"database" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

// Connect to PostgreSQL and check connection:
func (d Database) Connect() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     d.Address,
		User:     d.User,
		Password: d.Password,
		Database: d.Database,
	})

	if _, err := db.ExecOne("SELECT 1"); err != nil {
		return nil, err
	}

	return db, nil
}
