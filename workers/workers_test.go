package workers

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func init() {
	os.Setenv("LEVEL", "dev")
}

func testSimple(t *testing.T) {
	c, cancel := context.WithCancel(context.Background())
	w := New(nil)

	i := atomic.NewInt64(0)

	if err := w.Schedule(Options{
		Name:     "#1: 100 ms worker",
		Schedule: DelaySchedule(time.Millisecond * 100),
		Handler: func(ctx context.Context) error {
			i.Inc()
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(Options{
		Name:     "#2: 200 ms worker",
		Schedule: DelaySchedule(time.Millisecond * 200),
		Handler: func(ctx context.Context) error {
			i.Inc()
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(Options{
		Name:     "#3: 1 sec worker",
		Schedule: DelaySchedule(time.Second),
		Handler: func(ctx context.Context) error {
			panic("test")
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(Options{
		Name:     "test-1",
		Schedule: 0,
		Handler: func(ctx context.Context) error {
			panic("test")
		},
	}); !assert.Error(t, err) {
		t.Fatal("must be error")
	}

	if err := w.Schedule(Options{
		Name:     "test-2",
		Schedule: DelaySchedule(0),
		Handler: func(ctx context.Context) error {
			i.Inc()
			defer i.Dec()
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(Options{
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
	w.Start(c)

	<-c.Done()

	w.Wait(context.Background())

	assert.Equal(t, int64(6), i.Load())
}

func TestWorkers(t *testing.T) {
	t.Run("simple test workers", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Run("for-loop", func(t *testing.T) {
				testSimple(t)
			})
		}
	})

	t.Run("high way to hell", func(t *testing.T) {
		c, cancel := context.WithCancel(context.Background())
		w := New(nil)

		i := atomic.NewInt64(0)

		for n := 0; n < 1000; n++ {
			w.Schedule(Options{
				Name:     fmt.Sprintf("test-worker-%d", i),
				Schedule: DelaySchedule(time.Millisecond * 10),
				Handler: func(ctx context.Context) error {
					i.Inc()
					time.Sleep(time.Second * 2)
					defer i.Dec()
					return nil
				},
			})
		}

		w.Start(c)

		time.AfterFunc(time.Second, cancel)

		<-c.Done()

		w.Wait(context.Background())

		assert.Equal(t, int64(0), i.Load())
	})
}
