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

	fmt.Printf("[%s] Hello, workers!\n", time.Now().Format("15:04:05"))

	w := workers.New(ctx)

	// worker will run every 5 seconds
	// example of scheduler like time.Ticker
	if err := w.Schedule(&workers.ScheduleOptions{
		Name:     "every 5s",
		Schedule: time.Second * 5,
		Handler: func(ctx context.Context) error {
			fmt.Printf("[%s] worker #1 every 5 secs\n", time.Now().Format("15:04:05"))
			return nil
		},
	}); err != nil {
		panic(err)
	}

	w.Start()

	// worker will run every 13 seconds
	// example of scheduler like time.Ticker (using string parsing)
	step := atomic.NewUint32(0)

	if err := w.Schedule(&workers.ScheduleOptions{
		Name:     "every 13s",
		Schedule: time.Second * 13,
		Handler: func(ctx context.Context) error {
			fmt.Printf("[%s] worker #2 every 13 secs: STEP=%d\n", time.Now().Format("15:04:05"), step.Inc())
			return nil
		},
	}); err != nil {
		panic(err)
	}

	// worker will run every minutes at 12 secs
	// example of scheduler like UNIX cron
	// but with first element for seconds
	if err := w.Schedule(&workers.ScheduleOptions{
		Name:     "12 */1 * * * *",
		Schedule: time.Second * 13,
		Handler: func(ctx context.Context) error {
			fmt.Printf("[%s] worker #3 every minute at 12 secs\n", time.Now().Format("15:04:05"))
			return nil
		},
	}); err != nil {
		panic(err)
	}

	// worker will run as custom scheduler
	// example of scheduler (using workers.Schedule interface)
	delay := new(myDelayLock)
	if err := w.Schedule(&workers.ScheduleOptions{
		Name:     "every 5s",
		Schedule: time.Second * 13,
		Handler: func(ctx context.Context) error {
			if step.Load() > 4 && !delay.stop {
				fmt.Printf("[%s] worker #4: send command 'exit'\n", time.Now().Format("15:04:05"))
				delay.stop = true
				// delay canceling of context for 10 seconds
				time.AfterFunc(time.Second*10, cancel)
			}
			return nil
		},
	}); err != nil {
		panic(err)
	}

	// wait until context will be canceled
	<-ctx.Done()

	// stopping and wait workers
	w.Stop()

	fmt.Printf("[%s] All workers are stopped\n", time.Now().Format("15:04:05"))
}
