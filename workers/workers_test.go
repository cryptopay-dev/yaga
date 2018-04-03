package workers

import (
	"testing"
	"time"

	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestWorkers(t *testing.T) {
	log.Init()

	w := New()

	i := atomic.NewInt64(0)

	w.AddFunc("@every 1s", func() error {
		i.Inc()
		log.Infof("every 1s: %d", i.Load())
		return nil
	})

	w.AddFunc("@every 2s", func() error {
		i.Inc()
		log.Infof("every 2s: %d", i.Load())
		return nil
	})

	w.AddFunc("@every 1s", func() error {
		panic("test")
	})

	w.AddFunc("@every 6s", func() error {
		t.FailNow()
		return nil
	})

	w.Start()

	time.Sleep(time.Second * 5)

	w.Stop()

	assert.Equal(t, int64(7), i.Load())
}
