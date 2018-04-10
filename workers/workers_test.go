package workers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/cryptopay-dev/yaga/locker/local"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func testSimple(t *testing.T) {
	w := New(local.New(), nil, 100)
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
			panic("testing a logger of panic")
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(Options{
		Name:     "#4: 100 ms worker",
		Schedule: DelaySchedule(time.Millisecond * 100),
		TypeJob:  OnePerInstance,
		Handler: func(ctx context.Context) error {
			i.Inc()
			time.Sleep(time.Millisecond * 400)
			return errors.New("testing a logger of error")
		},
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

	c, cancel := context.WithTimeout(context.Background(), time.Millisecond*450)
	defer cancel()
	w.Start(c)

	assert.Equal(t, int64(7), i.Load())

	assert.Equal(t, 2, log.Count())
}

func TestWorkers(t *testing.T) {
	t.Run("multiple workers at one time", func(t *testing.T) {
		w := New(local.New(), nil, 10)
		log := newMockLogger()
		w.logger = log

		i := atomic.NewInt64(0)

		opts := Options{
			Name:     "my-best-test-worker",
			Schedule: DelaySchedule(time.Millisecond * 100),
			TypeJob:  OnePerInstance,
			Handler: func(ctx context.Context) error {
				time.Sleep(time.Millisecond * 400)
				i.Inc()
				return nil
			},
		}

		if err := w.Schedule(opts); !assert.NoError(t, err) {
			t.Fatal(err)
		}
		if err := w.Schedule(opts); !assert.Error(t, err) {
			t.Fatal("must be error")
		}

		c, cancel := context.WithTimeout(context.Background(), time.Millisecond*450)
		defer cancel()
		w.Start(c)

		assert.Equal(t, int64(1), i.Load())
		assert.Equal(t, 0, log.Count())
	})

	t.Run("simple test workers", func(t *testing.T) {
		testSimple(t)
	})

	t.Run("high way to hell", func(t *testing.T) {
		w := New(nil, nil, 100)
		log := newMockLogger()
		w.logger = log

		i := atomic.NewInt64(0)

		for n := 0; n < 210; n++ {
			w.Schedule(Options{
				Name:     fmt.Sprintf("test-worker-%d", n),
				Schedule: DelaySchedule(time.Millisecond * 90),
				Handler: func(ctx context.Context) error {
					i.Inc()
					time.Sleep(time.Millisecond * 20)
					defer i.Dec()
					return nil
				},
			})
		}

		c, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()
		w.Start(c)

		assert.Equal(t, int64(0), i.Load())
		assert.Equal(t, 10, log.Count())
	})
}
