# Yaga
Yaga is a project that provides a set of useful utilities for creating a golang service

![Documentation](./github/logo.png)

## Package description:

- [**Cacher**](./cacher) is a wrapper over Redis client with the basic methods for work that are described in the interface, with it you can write, get and delete data in Redis
- [**Cli**](./cli) with cli you can start the application with specifying the parameters. You can also use this package to run migrations or clean up the database
- [**Config**](./config) provides utilities for loading and validating config
- [**Conv**](./conv) provides functions for working with *decimal.Decimal*, converting from a string or float to decimal
- [**Decimal**](./decimal) is a wrapper over `github.com/shopspring/decimal` which allows you to work with decimal in a simplified way
- [**Errors**](./errors) is an error wrapper package that allows you to capture errors, convert them to a readable form and catch internal errors or panics, with conversion to a pretty logical error
- [**Locker**](./locker) is a wrapper over `github.com/bsm/redis-lock` for locks in Redis
- [**Logger**](./logger) provides the interface for its implementation for [zap](github.com/uber-go/zap) logger and for nop logger (dummy)
- [**Middlewares**](./middlewares) provides intermediate layers for authorizing and logging requests in web application
- [**Pprof**](./pprof) provides a utility for profiling with web interaction
- [**Report**](./report) helps to write metrics for data using `cryptopay-dev/go-metrics`
- [**Testdb**](./testdb) creates a connection to the test database to run tests
- [**Tracer**](./tracer) is a wrapper over the raven `github.com/getsentry/raven-go` client for the Sentry event/error logging system
- [**Web**](./web) allows you to run the web server using `github.com/labstack/echo` web framework with the necessary parameters
- [**Workers**](./workers) are tools to run goroutine and do some work on scheduling with a safe stop of their work

### Examples:

Some examples of packages you can find in [examples](./examples) folder or in `example_test.go` inside a particular package

### To run documentation based on the code:

Run in the console:
```
godoc -http=:6060
``` 
Then open `localhost:6060/`
