package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cryptopay-dev/yaga/workers"
	"go.uber.org/atomic"
)

func main() {
	w := workers.New()
	defer w.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("[%s] Hello, workers!\n", time.Now().Format("15:04:05"))

	// worker will run every 5 seconds
	w.AddJob(time.Second*5, func(context.Context) {
		fmt.Printf("[%s] worker #1 every 5 secs\n", time.Now().Format("15:04:05"))
	})

	// worker will run every 13 seconds
	step := atomic.NewUint32(0)
	w.AddJob(time.Second*13, func(context.Context) {
		fmt.Printf("[%s] worker #2 every 13 secs: STEP=%d\n", time.Now().Format("15:04:05"), step.Inc())
	})

	// worker will run every 1 seconds
	w.AddJob(time.Second, func(ctx context.Context) {
		if step.Load() == 4 {
			step.Inc()
			fmt.Printf("[%s] worker #3: send command 'exit'\n", time.Now().Format("15:04:05"))
			// delay canceling of context for 8 seconds
			time.AfterFunc(time.Second*8, cancel)
			return
		}
		select {
		case <-ctx.Done():
			fmt.Printf("[%s] worker #3: context is done\n", time.Now().Format("15:04:05"))
		default:
		}
	})

	// wait until context will be canceled
	<-ctx.Done()

	// stopping workers
	w.Stop()

	// wait until all workers will be stopped
	if err := w.Wait(nil); err != nil {
		fmt.Printf("[%s] Get error: %s\n", time.Now().Format("15:04:05"), err.Error())
	}

	fmt.Printf("[%s] All workers are stopped\n", time.Now().Format("15:04:05"))
}
