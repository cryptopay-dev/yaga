package commands

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDSN(t *testing.T) {
	var items = []struct {
		DSN   string
		Error error
	}{
		{
			DSN: "postgres://localhost:5432/database",
		},
		{
			DSN: "postgres://localhost:5432/database?sslmode=disable",
		},
		{
			DSN: "postgres://pg:pg@localhost:5432/database?sslmode=disable",
		},
		{
			DSN: "postgres://pg:pg@localhost:5432/database?sslmode=prefer",
		},
		{
			DSN: "postgres://pg:pg@localhost:5432/database?sslmode=allow",
		},
		{
			DSN:   "postgres://localhost:5432/database?sslmode=off",
			Error: errors.New("pg: sslmode 'off' is not supported"),
		},
		{
			DSN:   "mysql://localhost:5432/database?sslmode=off",
			Error: errors.New("pg: invalid scheme: mysql"),
		},
	}

	for _, item := range items {
		err := ParseDSN("database", item.DSN)

		if item.Error != nil {
			assert.EqualError(t, err, item.Error.Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
