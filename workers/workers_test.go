package workers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func testSimple(t *testing.T, iterN int) {
	c, cancel := context.WithCancel(context.Background())
	w := New(nil, nil, 10)
	log := newMockLogger()
	w.logger = log

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
		Name:     "#3: 400 ms worker",
		Schedule: DelaySchedule(time.Millisecond * 400),
		Handler: func(ctx context.Context) error {
			panic(fmt.Sprintf("(%d) testing a logger of panic", iterN))
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(Options{
		Name:     "#4: 100 ms worker",
		Schedule: DelaySchedule(time.Millisecond * 100),
		Handler: func(ctx context.Context) error {
			i.Inc()
			time.Sleep(time.Millisecond * 400)
			return fmt.Errorf("(%d) testing a logger of error", iterN)
		},
		TypeJob: OnePerInstance,
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(Options{
		Name:     "#5: error worker",
		Schedule: 0,
		Handler: func(ctx context.Context) error {
			t.Fatal("#5: must not be runned")
			return nil
		},
	}); !assert.Error(t, err) {
		t.Fatal("must be error")
	}

	if err := w.Schedule(Options{
		Name:     "#6: 1 ms worker",
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
		Name:     "#7: 2 sec worker",
		Schedule: DelaySchedule(time.Second * 2),
		Handler: func(ctx context.Context) error {
			t.Fatal("#7: must not be runned")
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	w.Start(c)
	time.AfterFunc(time.Millisecond*450, cancel)

	<-c.Done()

	assert.Equal(t, int64(10), i.Load())

	assert.Equal(t, 5, log.Count())
}

func TestWorkers(t *testing.T) {
	// TODO
	return
	t.Run("multiple workers at one time", func(t *testing.T) {
		c, cancel := context.WithCancel(context.Background())
		w := New(nil, nil, 10)

		i := atomic.NewInt64(0)

		opts := Options{
			Name:     "my-best-test-worker",
			Schedule: DelaySchedule(time.Millisecond * 900),
			TypeJob:  OnePerInstance,
			Handler: func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				i.Inc()
				return nil
			},
		}

		for n := 0; n < 10; n++ {
			w.Schedule(opts)
		}

		w.Start(c)

		time.AfterFunc(time.Second*2, cancel)

		<-c.Done()

		assert.Equal(t, int64(2), i.Load())
	})

	t.Run("simple test workers", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Run("for-loop", func(t *testing.T) {
				testSimple(t, i)
			})
		}
	})

	t.Run("high way to hell", func(t *testing.T) {
		c, cancel := context.WithCancel(context.Background())
		w := New(nil, nil, 10)

		i := atomic.NewInt64(0)

		for n := 0; n < 10; n++ {
			w.Schedule(Options{
				Name:     fmt.Sprintf("test-worker-%d", i),
				Schedule: DelaySchedule(time.Millisecond * 900),
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

		assert.Equal(t, int64(0), i.Load())
	})
}
