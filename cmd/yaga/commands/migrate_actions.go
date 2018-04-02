package commands

import (
	"fmt"
	"os"

	"github.com/cryptopay-dev/yaga/config"
	"github.com/cryptopay-dev/yaga/helpers/postgres"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/cryptopay-dev/yaga/migrate"
	"github.com/urfave/cli"
)

var stepFlag = cli.IntFlag{
	Name:  "steps",
	Value: 0,
	Usage: "steps count to up/down",
}

var dbFlag = cli.StringFlag{
	Name:  "db",
	Usage: "set database",
}

var dsnFlag = cli.StringFlag{
	Name:  "dsn",
	Value: "",
	Usage: "dsn (database source name)",
}

var mpathFlag = cli.StringFlag{
	Name:  "path",
	Usage: "migration path",
	Value: defaultMigratePath(),
}

func migrateFlags() []cli.Flag {
	return []cli.Flag{
		dbFlag,
		dsnFlag,
		stepFlag,
		mpathFlag,
	}
}

type migrateAct func(steps int) error
type migrateType uint

const (
	migrateUp = iota
	migrateDown
	migrateVersion
	migrateCleanup
	migrateList
	migratePlan
)

func (m migrateType) String() string {
	switch m {
	case migrateUp:
		return "Up"
	case migrateDown:
		return "Down"
	case migrateVersion:
		return "Version"
	case migrateCleanup:
		return "Cleanup"
	case migrateList:
		return "List"
	case migratePlan:
		return "Plan"
	default:
		return fmt.Sprintf("Unknown(%d)", m)
	}
}

func (m migrateType) needMigrations() bool {
	return m == migrateUp ||
		m == migrateDown ||
		m == migratePlan
}

func migrateAction(mtype migrateType) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) (err error) {
		if err = FetchDB(ctx, "database"); err != nil {
			log.Fatalf("can't find config file or dsn: %v", err)
		}

		if database := ctx.String("db"); len(database) != 0 {
			config.Set("database.database", database)
		}

		var steps = ctx.Int("steps")

		mpath := ctx.String("path")
		if mtype.needMigrations() {
			if _, err = os.Stat(mpath); err != nil {
				log.Fatalf("migration path not found: %v", err)
			}
		}

		pg, err := postgres.Connect("database")
		if err != nil {
			log.Fatalf("postgres connection error: %v", err)
		}

		m, err := migrate.New(migrate.Options{
			DB:   pg,
			Path: mpath,
		})

		if err != nil {
			log.Fatalf("migrate error: %v", err)
		}

		var action migrateAct

		switch mtype {
		case migrateUp:
			action = m.Up
		case migrateDown:
			action = m.Down
			if steps == 0 {
				steps = 1
			}
		case migrateVersion:
			action = func(int) error {
				version, errV := m.Version()
				if errV != nil {
					return err
				}
				log.Infof("database version: %d", version)
				return nil
			}
		case migrateList:
			action = func(int) error {
				items, errL := m.List()
				if errL != nil {
					return err
				}
				for _, item := range items {
					log.Infof(
						"%s -> %s",
						item.RealName(),
						item.CreatedAt,
					)
				}
				return nil
			}
		case migratePlan:
			action = func(int) error {
				items, errL := m.Plan()
				if errL != nil {
					return err
				}
				for _, item := range items {
					log.Infof("%s -> not applied", item.RealName())
				}
				return nil
			}
		default:
			log.Fatalf("migrate unknown action: %s", mtype)
		}

		if err = action(steps); err != nil {
			log.Fatalf("migrate action error: %v", err)
		}

		return
	}
}
