# Golang application metrics
> NATS application metrics package

[![Build Status](https://travis-ci.org/cryptopay-dev/go-metrics.svg?branch=master)](https://travis-ci.org/cryptopay-dev/go-metrics)
[![codecov](https://codecov.io/gh/cryptopay-dev/go-metrics/branch/master/graph/badge.svg)](https://codecov.io/gh/cryptopay-dev/go-metrics)
[![GoDoc](https://godoc.org/github.com/cryptopay-dev/go-metrics?status.svg)](https://godoc.org/github.com/cryptopay-dev/go-metrics)
[![Go Report Card](https://goreportcard.com/badge/github.com/cryptopay-dev/go-metrics)](https://goreportcard.com/report/github.com/cryptopay-dev/go-metrics)

## Installation
```bash
go get github.com/cryptopay-dev/go-metrics
```

## Default metrics tags
```
hostname - application host
app - application name
```

## Metrics
```json
    "alloc":         8810230, // memory allocated in bytes
    "alloc_objects": 123, // total heap objects allocated
    "gorotines":     10, // number of goroutines
    "gc":            1495532586, // timestamp of last GC
    "next_gc":       9000000, // heap size when GC will be run next time
    "pause_ns":      100 // pause time of GC
```

## Usage

## Basic
```go
package main

import (
    "log"

    "github.com/cryptopay.dev/go-metrics"
)

func main() {
    err := metrics.Setup("nats://localhost:4222", "application_name", "hostname")
    if err != nil {
        log.Fatal(err)
    }

    for i:=0; i<10; i++ {
        // You metrics will be reported at application_name:metric
        // You metrics will be send to: mymetric,user=test@example.com counter=1,gauge=true,string=name
        // E.t.c.
        err = metrics.SendAndWait("metric", metrics.M{
            "counter": i,
            "gauge": true,
            "string": "name",
        }, metrics.T{
            "user": "test@example.com"
        }, "mymetric")

        if err != nil {
            log.Fatal(err)
        }
    }
}
```

## More stuff
You can find more examples in `/examples` folder 