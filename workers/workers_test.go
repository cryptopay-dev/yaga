package workers

import (
	"context"
	"testing"
	"time"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestWorkers(t *testing.T) {
	log.Init()

	c, cancel := context.WithCancel(context.Background())

	w := New(c)

	i := atomic.NewInt64(0)

	if err := w.Schedule(&ScheduleOptions{
		Duration: time.Second,
		Handler: func(ctx context.Context) error {
			i.Inc()
			log.Infof("every 100ms: %d", i.Load())
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(&ScheduleOptions{
		Duration: time.Second * 2,
		Handler: func(ctx context.Context) error {
			i.Inc()
			log.Infof("every 200ms: %d", i.Load())
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(&ScheduleOptions{
		Duration: time.Second,
		Handler: func(ctx context.Context) error {
			panic("test")
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := w.Schedule(&ScheduleOptions{
		Duration: time.Second * 10,
		Handler: func(ctx context.Context) error {
			t.Fatal("must not be runned")
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}

	w.Start()

	time.AfterFunc(time.Second*5, cancel)

	<-c.Done()

	w.Stop()

	assert.Equal(t, int64(7), i.Load())
}
