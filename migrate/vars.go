package migrate

import (
	"errors"
	"fmt"
)

const (
	errFileNamingTpl      = "bad file name '%s', must be like '<timestamp>_something.<up|down>.sql'"
	errFileVersionTpl     = "bad file version '%s', must be greater than 0"
	errVersionNotEqualTpl = "version of 'up' and 'down' migrations must be equal: %d != %d"

	fileNameTpl = "%d_%s.%s.sql"

	sqlCreateSchema = `CREATE SCHEMA IF NOT EXISTS ?`
	sqlNewVersion   = `INSERT INTO ? (version, created_at) VALUES (?, now())`
	sqlRemVersion   = `DELETE FROM ? WHERE version = ?`
	sqlGetVersion   = `SELECT version FROM ? ORDER BY id DESC LIMIT 1`
	sqlCreateTable  = `
CREATE TABLE IF NOT EXISTS ? (
	id serial,
	version bigint UNIQUE,
	created_at timestamptz,
	PRIMARY KEY(id)
)`
)

var (
	schemaName string
	tableName  = "migrations"

	// ErrNoDB set to Options
	ErrNoDB = fmt.Errorf("no db")
	// ErrNoLogger set to Options
	ErrNoLogger = fmt.Errorf("no logger")
	// ErrDirNotExists when migration path not exists
	ErrDirNotExists = fmt.Errorf("migrations dir not exists")
	// ErrBothMigrateTypes when up or down migration file not found
	ErrBothMigrateTypes = errors.New("migration must have up and down files")
)
