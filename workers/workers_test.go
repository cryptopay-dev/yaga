package workers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func testSimple(t *testing.T) {
	t.Parallel()

	c, cancel := context.WithCancel(context.Background())

	w := New(c)

	//start1 := time.Now()
	//start2 := time.Now()

	i := atomic.NewInt64(0)

	if err := w.Schedule(&ScheduleOptions{
		Name:     "#1: 100 ms worker",
		Schedule: DelaySchedule(time.Millisecond * 100),
		Handler: func(ctx context.Context) error {
			i.Inc()
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(&ScheduleOptions{
		Name:     "#2: 200 ms worker",
		Schedule: DelaySchedule(time.Millisecond * 200),
		Handler: func(ctx context.Context) error {
			i.Inc()
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(&ScheduleOptions{
		Name:     "#3: 1 sec worker",
		Schedule: DelaySchedule(time.Second),
		Handler: func(ctx context.Context) error {
			panic("test")
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(&ScheduleOptions{
		Name:     "#4: 10 sec worker",
		Schedule: DelaySchedule(time.Second * 10),
		Handler: func(ctx context.Context) error {
			t.Fatal("must not be runned")
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	time.AfterFunc(time.Millisecond*450, cancel)
	w.Start()

	<-c.Done()

	w.Stop()

	assert.Equal(t, int64(6), i.Load())
}

func TestWorkers(t *testing.T) {
	log.Init()

	t.Run("simple test workers", func(t *testing.T) {
		t.Parallel()

		for i := 0; i < 10; i++ {
			t.Run("for-loop", func(t *testing.T) {
				testSimple(t)
			})
		}
	})

	t.Run("high way to hell", func(t *testing.T) {
		t.Parallel()

		c, cancel := context.WithCancel(context.Background())
		w := New(c)

		i := atomic.NewInt64(0)

		for n := 0; n < 1000; n++ {
			w.Schedule(&ScheduleOptions{
				Name:     fmt.Sprintf("test-worker-%d", i),
				Schedule: time.Millisecond,
				Handler: func(ctx context.Context) error {
					i.Inc()
					time.Sleep(time.Second * 2)
					defer i.Dec()
					return nil
				},
			})
		}

		w.Start()

		time.AfterFunc(time.Second, cancel)

		<-c.Done()

		w.Stop()

		assert.Equal(t, int64(0), i.Load())
	})
}
