package cli

import (
	"fmt"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/mattes/migrate"

	// Migration libs
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

// MigrationLogger abstraction for migrate-logger
type MigrationLogger struct {
	logger logger.Logger
}

// NewMigrationLogger creates abstraction for migrate-logger
func NewMigrationLogger(logger logger.Logger) MigrationLogger {
	return MigrationLogger{
		logger: logger,
	}
}

// Printf implementation of migrate.Logger
func (m MigrationLogger) Printf(format string, v ...interface{}) {
	m.logger.Info(fmt.Sprintf(format, v...))
}

// Verbose implementation of migrate.Logger
func (m MigrationLogger) Verbose() bool {
	return true
}

// MigrationDirection type
type MigrationDirection int

const (
	_ MigrationDirection = iota

	// MigrationUp to migrate step up
	MigrationUp
	// MigrationDown to migrate step down
	MigrationDown
)

// MigrateOptions struct used for Migrate method
type MigrateOptions struct {
	Steps     int
	Direction MigrationDirection
}

// MigrateOption closure
type MigrateOption func(o *MigrateOptions)

//  newMigrateOptions converts slice of closures to MigrateOptions-field
func newMigrateOptions(opts ...MigrateOption) MigrateOptions {
	// Defaults:
	var options = MigrateOptions{
		Steps: 1,
	}

	// Set options:
	for _, o := range opts {
		o(&options)
	}

	return options
}

// MigrateDirection closure to set field in MigrateOptions
func MigrateDirection(direction MigrationDirection) MigrateOption {
	return func(o *MigrateOptions) {
		o.Direction = direction
	}
}

// MigrateSteps closure to set field in MigrateOptions
func MigrateSteps(steps int) MigrateOption {
	return func(o *MigrateOptions) {
		o.Steps = steps
	}
}

// Migrate db using MigrateOption's
func (o Options) Migrate(opts ...MigrateOption) error {
	var options = newMigrateOptions(opts...)

	m, err := o.migrate()
	if err != nil {
		return err
	}

	defer func() {
		var (
			version uint
			dirty   bool
		)

		if version, dirty, err = m.Version(); err != nil {
			m.Log.Printf("Error while retrieve version: %v", err)
			return
		}

		m.Log.Printf("Version %v, Dirty %v", version, dirty)
	}()

	var migrationErr error
	switch options.Direction {
	case MigrationUp:
		migrationErr = m.Up()
	case MigrationDown:
		migrationErr = m.Steps(-1 * options.Steps)
	}

	if migrationErr != nil {
		if migrationErr == migrate.ErrNoChange {
			return nil
		}

		return err
	}

	return nil
}

// migrate returns a new Migrate instance from a source URL and a database URL.
// The URL scheme is defined by each driver
func (o Options) migrate() (*migrate.Migrate, error) {
	uri := "postgres://%s:%s@%s/%s?sslmode=disable"

	m, err := migrate.New(
		"file://./migrations",
		fmt.Sprintf(
			uri,
			o.DB.Options().User,
			o.DB.Options().Password,
			o.DB.Options().Addr,
			o.DB.Options().Database,
		))
	if err != nil {
		return nil, err
	}

	m.Log = NewMigrationLogger(o.Logger)
	return m, nil
}
