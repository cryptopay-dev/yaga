package migrate

import (
	"errors"
	"fmt"
	"testing"

	"github.com/cryptopay-dev/yaga/helpers/testdb"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/stretchr/testify/assert"
)

var (
	// Nop:
	defaultLogger = nop.New()
	// Debug:
	//defaultLogger = zap.New(zap.Development)
	errEmpty = errors.New("[empty error]")
)

const (
	zero int64 = 0
)

type mockDB struct {
	*pg.DB
	*pg.Tx
}

func (m *mockDB) RunInTransaction(fn func(*pg.Tx) error) error {
	if m.DB == nil {
		return errEmpty
	}
	return m.DB.RunInTransaction(fn)
}
func (m *mockDB) Exec(query interface{}, params ...interface{}) (orm.Result, error) {
	if m.Tx == nil {
		return nil, errEmpty
	}
	return m.Tx.Exec(query, params...)
}
func (m *mockDB) ExecOne(query interface{}, params ...interface{}) (orm.Result, error) {
	if m.Tx == nil {
		return nil, errEmpty
	}
	return m.Tx.ExecOne(query, params...)
}
func (m *mockDB) Query(model, query interface{}, params ...interface{}) (orm.Result, error) {
	if m.Tx == nil {
		return nil, errEmpty
	}
	return m.Tx.Query(model, query, params...)
}
func (m *mockDB) QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error) {
	if m.Tx == nil {
		return nil, errEmpty
	}
	return m.Tx.QueryOne(model, query, params...)
}

func TestUpDown(t *testing.T) {
	var db = testdb.GetTestDB().DB

	t.Run("Good", func(t *testing.T) {
		m, errNew := New(Options{
			DB:     db,
			Path:   "./fixtures/good",
			Logger: defaultLogger,
		})

		if !assert.NoError(t, errNew) {
			t.FailNow()
		}

		if errUp := m.Up(0); !assert.NoError(t, errUp) {
			t.FailNow()
		}

		if errDown := m.Down(0); !assert.NoError(t, errDown) {
			t.FailNow()
		}
	})

	t.Run("Bad", func(t *testing.T) {
		m, errNew := New(Options{
			DB:     db,
			Path:   "./fixtures/bad",
			Logger: defaultLogger,
		})

		if !assert.NoError(t, errNew) {
			t.FailNow()
		}

		if errUp := m.Up(0); !assert.Error(t, errUp) {
			t.FailNow()
		}

		mig, ok := m.(*migrate)
		if !assert.True(t, ok) {
			t.FailNow()
		}

		db.RunInTransaction(func(tx *pg.Tx) error {
			for _, item := range mig.migrations {
				if errVer := addVersion(tx, item.Version); !assert.NoError(t, errVer) {
					return errVer
				}
			}

			return nil
		})

		if errDown := m.Down(0); !assert.Error(t, errDown) {
			t.FailNow()
		}

		db.RunInTransaction(func(tx *pg.Tx) error {
			for _, item := range mig.migrations {
				if errVer := remVersion(tx, item.Version); !assert.NoError(t, errVer) {
					return errVer
				}
			}

			return nil
		})
	})
}

func TestNew(t *testing.T) {
	var (
		err error
		db  = testdb.GetTestDB().DB
	)

	t.Run("Good case", func(t *testing.T) {
		if err = db.RunInTransaction(func(tx *pg.Tx) error {
			tx.Exec(`TRUNCATE ?`, getTableName())
			m, errNew := New(Options{
				DB:     &mockDB{DB: db, Tx: tx},
				Path:   "./fixtures/good",
				Logger: defaultLogger,
			})

			if !assert.NoError(t, errNew) {
				return errNew
			}

			if !assert.NotNil(t, m) {
				return errors.New("empty migrator")
			}

			if ver, errVer := m.Version(); !assert.NoError(t, errVer) {
				return errVer
			} else if !assert.Equal(t, zero, ver) {
				return fmt.Errorf("wrong migration version: %d != %d", zero, ver)
			}

			var i int64

			for i = 1; i <= 10; i++ {
				if _, errVer := tx.Exec(sqlNewVersion, getTableName(), i); errVer != nil {
					return errVer
				}

				if ver, errVer := m.Version(); !assert.NoError(t, errVer) {
					return errVer
				} else if !assert.Equal(t, i, ver) {
					return fmt.Errorf("wrong migration version: %d != %d", i, ver)
				}
			}

			// to reject inserts
			return errEmpty
		}); err != nil && err != errEmpty {
			t.Fatal(err)
		}
	})

	t.Run("Bad case #1", func(t *testing.T) {
		_, err = New(Options{
			DB: nil,
		})
		assert.Error(t, err)
	})

	t.Run("Bad case #2", func(t *testing.T) {
		_, err = New(Options{
			DB:   db,
			Path: "",
		})
		assert.Error(t, err)
	})

	t.Run("Bad case #3", func(t *testing.T) {
		_, err = New(Options{
			DB:   db,
			Path: "/no/such/dir",
		})
		assert.Error(t, err)
	})

	t.Run("Bad case #4", func(t *testing.T) {
		_, err = New(Options{
			DB:   db,
			Path: "/dev/null",
		})
		assert.Error(t, err)
	})

	t.Run("Bad case #5", func(t *testing.T) {
		_, err = New(Options{
			DB:   db,
			Path: "./fixtures/bad",
		})
		assert.Error(t, err)
	})

	t.Run("Bad case #6", func(t *testing.T) {
		_, err = New(Options{
			DB:     db,
			Path:   "./fixtures/good",
			Logger: nil,
		})
		assert.Error(t, err)
	})
}

func TestMigrate_Version(t *testing.T) {
	var db = testdb.GetTestDB().DB

	db.RunInTransaction(func(tx *pg.Tx) error {
		tx.Exec(`TRUNCATE ?`, getTableName())
		m := migrate{
			Options: Options{
				DB:     &mockDB{DB: db, Tx: tx},
				Path:   "./fixtures/good",
				Logger: defaultLogger,
			},
		}

		// Good case:
		if ver, errVer := m.Version(); !assert.NoError(t, errVer) {
			return errVer
		} else if !assert.Equal(t, zero, ver) {
			return fmt.Errorf("wrong migration version: %d != %d", zero, ver)
		}

		// Bad case #1:
		m.DB = &mockDB{Tx: nil}

		if _, errVer := m.Version(); !assert.Error(t, errVer) {
			return errVer
		}

		// Bad case #2:
		m.DB = &mockDB{Tx: nil}

		if _, errVer := m.Version(); !assert.Error(t, errVer) {
			return errVer
		}

		// to reject inserts
		return errEmpty
	})
}
