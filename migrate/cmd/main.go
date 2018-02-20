package main

import (
	"flag"

	"github.com/cryptopay-dev/yaga/logger/zap"
	"github.com/cryptopay-dev/yaga/migrate"
	"github.com/go-pg/pg"
)

const (
	// Up action
	Up = "up"
	// Down action
	Down = "down"
	// List action
	List = "list"
)

// Options for migrate cmd
type Options struct {
	Addr     string
	User     string
	Database string
	Password string
	Path     string
	Type     string
	Steps    int
}

// Validate migrate options
func (o Options) Validate() bool {
	return len(o.Addr) > 0 &&
		len(o.User) > 0 &&
		len(o.Database) > 0 &&
		len(o.Path) > 0 &&
		(o.Type == Up || o.Type == Down || o.Type == List)
}

// PG options
func (o Options) PG() *pg.Options {
	return &pg.Options{
		Addr:     o.Addr,
		User:     o.User,
		Database: o.Database,
		Password: o.Password,
	}
}

var (
	log  = zap.New(zap.Development)
	opts = new(Options)
)

func main() {
	flag.StringVar(&opts.Password, "p", "", "password")
	flag.StringVar(&opts.Addr, "a", "localhost:5432", "address")
	flag.StringVar(&opts.Database, "d", "", "db name")
	flag.StringVar(&opts.User, "u", "", "username")
	flag.StringVar(&opts.Path, "src", "./migrations", "path to migrations")
	flag.StringVar(&opts.Type, "t", "", "migrate action up/down")
	flag.IntVar(&opts.Steps, "s", 0, "steps")

	flag.Parse()

	if !opts.Validate() {
		flag.PrintDefaults()
		return
	}

	m, err := migrate.New(migrate.Options{
		DB:     pg.Connect(opts.PG()),
		Path:   opts.Path,
		Logger: log,
	})

	if err != nil {
		log.Fatal(err)
	}

	switch opts.Type {
	case Up:
		if err = m.Up(opts.Steps); err != nil {
			log.Fatal(err)
		}
	case Down:
		if err = m.Down(opts.Steps); err != nil {
			log.Fatal(err)
		}
	case List:
		items, err := m.List()
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range items {
			log.Infof("%d %s -> %s", item.Version, item.Name, item.CreatedAt)
		}
	default:
		log.Fatal("unknown migrate action")
	}

	log.Info("all done")
}
