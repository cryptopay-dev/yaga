# Yaga
Yaga is a project that provides a set of useful utilities for creating a golang service

![Documentation](./github/logo.png)

## Package description:

- [**Cacher**](./cacher) is a wrapper over Redis client with the basic methods for work that are described in the interface, with it you can write, get and delete data in Redis
- [**Cli**](./cli) with cli you can start the application with specifying the parameters. You can also use this package to run migrations or clean up the database
- [**Collection**](./helpers/collection) provides simple method to list your data from database
- [**Config**](./config) provides utilities for loading and validating config
- [**Conv**](./conv) provides functions for working with *decimal.Decimal*, converting from a string or float to decimal
- [**Decimal**](./decimal) is a wrapper over `github.com/shopspring/decimal` which allows you to work with decimal in a simplified way
- [**Doc**](./doc) tool for echo, allowing you to make documentation based on the swagger file and `rebilly.github.io/ReDoc`
- [**Errors**](./errors) is an error wrapper package that allows you to capture errors, convert them to a readable form and catch internal errors or panics, with conversion to a pretty logical error
- [**Locker**](./locker) is a wrapper over `github.com/bsm/redis-lock` for locks in Redis
- [**Logger**](./logger) provides the interface for its implementation for [zap](github.com/uber-go/zap) logger and for nop logger (dummy)
- [**Mail**](./mail) service for send emails
- [**Middlewares**](./middlewares) provides intermediate layers for authorizing and logging requests in web application
- [**Migrator**](./migrate) this package allows you to run migrations on your PostgreSQL database
- [**Pprof**](./pprof) provides a utility for profiling with web interaction
- [**Report**](./report) helps to write metrics for data using `cryptopay-dev/go-metrics`
- [**Testdb**](./helpers/testdb) creates a connection to the test database to run tests
- [**Tracer**](./tracer) is a wrapper over the raven `github.com/getsentry/raven-go` client for the Sentry event/error logging system
- [**Web**](./web) allows you to run the web server using `github.com/labstack/echo` web framework with the necessary parameters
- [**Workers**](./workers) are tools to run goroutine and do some work on scheduling with a safe stop of their work

## Yaga commandline tool:

### Install:

```
go get github.com/cryptopay-dev/yaga/cmd/...
```

### Usage:
```
$ go get github.com/cryptopay-dev/yaga/cmd/...
$ yaga
NAME:
   yaga - Yaga command line tool

USAGE:
   yaga [global options] command [command options] [arguments...]

VERSION:
   v1.9.11 (2018-03-07 11:55:53 +0300)

COMMANDS:
     new, n   new <work-dir>
     help, h  Shows a list of commands or help for one command
   Migrate commands:
     migrate:create, m:c    new <migration-name> --path=<to-migrations>
     migrate:up, m:u        up --steps=<count> --dsn=<DSN> --db=<db-name> --path=<to-migrations>
     migrate:down, m:d      down --steps=<count> --dsn=<DSN> --db=<db-name> --path=<to-migrations>
     migrate:version, m:v   version --db=<db-name> --dsn=<DSN>
     migrate:list, m:l      list --db=<db-name> --dsn=<DSN>
     migrate:plan, m:p      plan --db=<db-name> --dsn=<DSN> --db=<db-name> --path=<to-migrations>
     migrate:cleanup, m:cl  cleanup --db=<db-name> --dsn=<DSN>

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Examples:

Some examples of packages you can find in [examples](./examples) folder or in `example_test.go` inside a particular package

### To run documentation based on the code:

Run in the console:
```
godoc -http=:6060
``` 
Then open `localhost:6060/`
