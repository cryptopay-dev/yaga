package migrate

import (
	"os"
	"sort"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Options for migrator
type Options struct {
	// DB connection
	DB DB
	// Path to migrations files
	Path string
	// Logger
	Logger logger.Logger
}

// Migrator interface
type Migrator interface {
	Up(steps int) error
	Down(steps int) error
	Version() (int64, error)
}

// DB interface
type DB interface {
	RunInTransaction(fn func(*pg.Tx) error) error
	Exec(query interface{}, params ...interface{}) (orm.Result, error)
	ExecOne(query interface{}, params ...interface{}) (orm.Result, error)
	Query(model, query interface{}, params ...interface{}) (orm.Result, error)
	QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error)
}

// migrate is implementation of Migrator
type migrate struct {
	Options
	migrations
}

// migration item
type migration struct {
	Version int64
	Name    string
	Up      func(DB) error
	Down    func(DB) error
}

// migrations slice
type migrations []*migration

// New creates new Migrator
func New(opts Options) (Migrator, error) {
	var (
		err   error
		items migrations
		files []os.FileInfo
	)

	if opts.DB == nil {
		return nil, ErrNoDB
	}

	if opts.Logger == nil {
		return nil, ErrNoLogger
	}

	if files, err = findMigrations(opts.Path); err != nil {
		return nil, err
	}

	if items, err = extractMigrations(opts.Logger, opts.Path, files); err != nil {
		return nil, err
	}

	return &migrate{
		Options:    opts,
		migrations: items,
	}, nil
}

// createTables for migrations
func createTables(db DB) error {
	var err error
	if len(schemaName) > 0 {
		if _, err = db.Exec(
			sqlCreateSchema,
			pg.Q(schemaName),
		); err != nil {
			return err
		}
	}

	_, err = db.Exec(sqlCreateTable, pg.Q(tableName))

	return err
}

// Up, roll up multiple migrations
func (m *migrate) Up(steps int) error {
	var (
		err     error
		version int64
		count   = len(m.migrations)
	)

	if steps < 0 {
		return nil
	}

	if steps == 0 {
		steps = count
	}

	if version, err = m.Version(); err != nil {
		return err
	}

	items := make(migrations, count)

	copy(items, m.migrations)

	sort.Slice(items, func(i, j int) bool {
		return items[i].Version < items[j].Version
	})

	for i, item := range items {
		if steps <= 0 {
			break
		}

		if item.Version <= version {
			continue
		}

		m.Logger.Infof("(%d) migrate up to: %d_%s", i+1, item.Version, item.Name)
		if err = item.Up(m.DB); err != nil {
			return err
		}

		steps--
	}

	return nil
}

// Down rollback some migrations
func (m *migrate) Down(steps int) error {
	var (
		err     error
		version int64
		count   = len(m.migrations)
	)

	if steps < 0 {
		return nil
	}

	if steps > count || steps == 0 {
		steps = count
	}

	if version, err = m.Version(); err != nil {
		return err
	}

	if version <= 0 {
		return nil
	}

	items := make(migrations, count)

	copy(items, m.migrations)

	sort.Slice(items, func(i, j int) bool {
		return items[i].Version > items[j].Version
	})

	for _, item := range items {
		if steps <= 0 {
			break
		}

		if item.Version > version {
			continue
		}

		m.Logger.Infof("(%d) migrate down to: %d_%s", steps, item.Version, item.Name)
		if err = item.Down(m.DB); err != nil {
			return err
		}

		steps--
	}

	return nil
}

// Version fetching from database
func (m *migrate) Version() (version int64, err error) {
	version = -1

	if err = createTables(m.DB); err != nil {
		return
	}

	if _, err = m.DB.QueryOne(
		pg.Scan(&version),
		sqlGetVersion,
		getTableName(),
	); err != nil && err == pg.ErrNoRows {
		err = nil
		version = 0
	}

	return
}
