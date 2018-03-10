package internal

import (
	"errors"
	"testing"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/stretchr/testify/assert"
)

type testConfigWithDB struct {
	Database config.Database
}

var (
	testDB     = &config.Database{Address: "localhost:5432", Database: "database", User: "pg", Password: "pg"}
	badConfig  = struct{ test string }{}
	goodConfig = testConfigWithDB{Database: *testDB}
)

func TestParseDSN(t *testing.T) {
	var items = []struct {
		DSN     string
		Options *config.Database
		Error   error
	}{
		{
			DSN:     "postgres://localhost:5432/database",
			Options: &config.Database{Address: "localhost:5432", Database: "database", User: "postgres"},
		},
		{
			DSN:     "postgres://localhost:5432/database?sslmode=disable",
			Options: &config.Database{Address: "localhost:5432", Database: "database", User: "postgres"},
		},
		{
			DSN:     "postgres://pg:pg@localhost:5432/database?sslmode=disable",
			Options: testDB,
		},
		{
			DSN:     "postgres://pg:pg@localhost:5432/database?sslmode=prefer",
			Options: testDB,
		},
		{
			DSN:     "postgres://pg:pg@localhost:5432/database?sslmode=allow",
			Options: testDB,
		},
		{
			DSN:     "postgres://localhost:5432/database?sslmode=off",
			Options: nil,
			Error:   errors.New("pg: sslmode 'off' is not supported"),
		},
		{
			DSN:     "mysql://localhost:5432/database?sslmode=off",
			Options: nil,
			Error:   errors.New("pg: invalid scheme: mysql"),
		},
	}

	for _, item := range items {
		opt, err := ParseDSN(item.DSN)

		if item.Error != nil {
			assert.Nil(t, opt)
			assert.EqualError(t, err, item.Error.Error())
		} else {
			assert.Equal(t, item.Options, opt)
			assert.NoError(t, err)
		}
	}
}

func TestParseConfig(t *testing.T) {
	var items = []struct {
		Config  interface{}
		Options *config.Database
		Error   error
	}{
		{
			Config:  goodConfig,
			Options: testDB,
		},
		{
			Config:  &goodConfig,
			Options: testDB,
		},
		{
			Config:  badConfig,
			Options: nil,
			Error:   errors.New("config hasn't database param"),
		},
		{
			Config:  &badConfig,
			Options: nil,
			Error:   errors.New("config hasn't database param"),
		},
	}

	for _, item := range items {
		opt, err := ParseConfig(item.Config)

		if item.Error != nil {
			assert.Nil(t, opt)
			assert.EqualError(t, err, item.Error.Error())
		} else {
			assert.Equal(t, item.Options, opt)
			assert.NoError(t, err)
		}
	}
}
