package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cryptopay-dev/yaga/workers"
	"go.uber.org/atomic"
)

type myDelayLock struct {
	stop bool
}

func (m *myDelayLock) Next(t time.Time) time.Time {
	if m.stop {
		fmt.Printf("[%s] instance will be stopped\n", time.Now().Format("15:04:05"))
		// we stop the worker using zero time
		return time.Time{}
	}

	// we will plan on every second
	return t.Add(time.Second)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := workers.New(ctx)

	fmt.Printf("[%s] Hello, workers!\n", time.Now().Format("15:04:05"))

	// worker will run every 5 seconds
	// example of scheduler like time.Ticker
	err := w.Schedule(workers.Options{
		Name:     "worker #1",
		Schedule: time.Second * 5,
		Handler: func(context.Context) error {
			fmt.Printf("[%s] worker #1 every 5 secs\n", time.Now().Format("15:04:05"))
			return nil
		},
	})
	if err != nil {
		panic(err)
	}

	w.Start()

	// worker will run every 13 seconds
	// example of scheduler like time.Ticker (using string parsing)
	step := atomic.NewUint32(0)
	err = w.Schedule(workers.Options{
		Name:     "worker #2",
		Schedule: "@every 13s",
		Handler: func(context.Context) error {
			fmt.Printf("[%s] worker #2 every 13 secs: STEP=%d\n", time.Now().Format("15:04:05"), step.Inc())
			return nil
		},
	})
	if err != nil {
		panic(err)
	}

	// worker will run every minutes at 12 secs
	// example of scheduler like UNIX cron
	// but with first element for seconds
	err = w.Schedule(workers.Options{
		Name:     "worker #3",
		Schedule: "12 */1 * * * *",
		Handler: func(context.Context) error {
			fmt.Printf("[%s] worker #3 every minute at 12 secs\n", time.Now().Format("15:04:05"))
			return nil
		},
	})
	if err != nil {
		panic(err)
	}

	// worker will run as custom scheduler
	// example of scheduler (using workers.Schedule interface)
	delay := new(myDelayLock)
	err = w.Schedule(workers.Options{
		Name:     "worker #4",
		Schedule: delay,
		Handler: func(context.Context) error {
			if step.Load() > 4 && !delay.stop {
				fmt.Printf("[%s] worker #4: send command 'exit'\n", time.Now().Format("15:04:05"))
				delay.stop = true
				// delay canceling of context for 10 seconds
				time.AfterFunc(time.Second*10, cancel)
			}
			return nil
		},
	})
	if err != nil {
		panic(err)
	}

	// wait until context will be canceled
	<-ctx.Done()

	fmt.Printf("[%s] workers are stopping\n", time.Now().Format("15:04:05"))

	// stopping workers
	w.Stop()

	// wait until all workers will be stopped
	w.Wait(context.Background())

	fmt.Printf("[%s] All workers are stopped\n", time.Now().Format("15:04:05"))
}
