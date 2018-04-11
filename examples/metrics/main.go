package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/cryptopay-dev/yaga/metrics"
	"github.com/cryptopay-dev/yaga/metrics/prometheus"
)

func main() {
	addr := ":18080"

	log.Info("Create prometheus metrics provider")
	provider := prometheus.NewProvider(addr)

	log.Info("Setup metrics provider")
	metrics.SetProvider(provider)

	log.Infof("Starting prometheus web server with address: %s", addr)
	provider.StartWebServer()

	log.Info("Write `example_counter` counter, increase to 5")
	metrics.IncrementCounter("example_counter", 5)

	log.Info("Write `example_gauage` gauage, increase to 5")
	metrics.IncrementGauge("example_gauage", 5)

	log.Info("Write `example_gauage` gauage, decrease to 3")
	metrics.IncrementGauge("example_gauage", -2)

	log.Info("Write `example_summary` summary, write 1s value")
	metrics.Observe("example_summary", time.Second.Seconds())

	log.Infof("Open %s and check metrics values", addr)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-ch

	log.Info("Gotten stop signal")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Info("Shutdown prometheus web server with 1s timeout")
	err := provider.StopWebServer(ctx)
	log.Info("Stop prometheus web server return error: %s", err)
}
