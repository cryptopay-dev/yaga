package migrate

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cryptopay-dev/yaga/logger"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/types"
)

// getTableName for quote table name
func getTableName() types.ValueAppender {
	return pg.Q(tableName)
}

// extractAttributes, such as  version, name and migration-type
func extractAttributes(filename string) (version int64, name, mType string, err error) {
	parts := strings.SplitN(filename, "_", 2)

	if len(parts) != 2 {
		err = fmt.Errorf(errFileNamingTpl, filename)
		return
	}

	if version, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
		return
	} else if version <= 0 {
		err = fmt.Errorf(errFileVersionTpl, parts[0])
		return
	}

	parts = strings.SplitN(parts[1], ".", 3)

	if len(parts) != 3 || parts[1] != "down" && parts[1] != "up" {
		err = fmt.Errorf(errFileNamingTpl, filename)
		return
	}

	name, mType = parts[0], parts[1]

	return
}

// findMigrations in specified folder (path)
func findMigrations(path string) ([]os.FileInfo, error) {
	var (
		err error
		dir *os.File
	)

	if dir, err = os.Open(path); os.IsNotExist(err) {
		return nil, ErrDirNotExists
	}

	return dir.Readdir(0)
}

// updateVersion abstraction
type updateVersion func(tx *pg.Tx, version int64) error

// remVersion migration from database
func remVersion(tx *pg.Tx, version int64) error {
	_, err := tx.Exec(sqlRemVersion, getTableName(), version)
	return err
}

// addVersion migration to database
func addVersion(tx *pg.Tx, version int64) error {
	_, err := tx.Exec(sqlNewVersion, getTableName(), version)
	return err
}

// doMigrate closure
func doMigrate(version int64, sql string, fn updateVersion) func(db DB) error {
	return func(db DB) error {
		return db.RunInTransaction(func(tx *pg.Tx) error {
			if _, errQuery := tx.Exec(sql); errQuery != nil {
				return errQuery
			}

			if errVersion := fn(tx, version); errVersion != nil {
				return errVersion
			}

			return nil
		})
	}
}

// extractMigrations, find files in migration folder and convert to migration-item
func extractMigrations(log logger.Logger, path string, files []os.FileInfo) (migrations, error) {
	var (
		err          error
		data         []byte
		migrateParts = make(map[string]*migration)
		items        migrations
	)

	for _, file := range files {
		log.Infof("Prepare migration file: %s", file.Name())

		if data, err = ioutil.ReadFile(path + "/" + file.Name()); err != nil {
			return nil, err
		}

		ver, name, mType, err := extractAttributes(file.Name())
		if err != nil {
			return nil, err
		}

		m, ok := migrateParts[name]
		if !ok {
			m = &migration{
				Version: ver,
				Name:    name,
			}
		} else if m.Version != ver {
			return nil, fmt.Errorf(errVersionNotEqualTpl, m.Version, ver)
		}

		switch mType {
		case "up":
			m.Up = doMigrate(ver, string(data), addVersion)
		case "down":
			m.Down = doMigrate(ver, string(data), remVersion)
		}

		migrateParts[name] = m
	}

	items = make(migrations, 0, len(migrateParts))
	for name, m := range migrateParts {
		log.Infof("Prepare migration: %s", name)

		if m.Down == nil || m.Up == nil {
			return nil, ErrBothMigrateTypes
		}

		items = append(items, m)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Version < items[j].Version
	})

	return items, nil
}
