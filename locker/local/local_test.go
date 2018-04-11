package local

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestLock(t *testing.T) {
	t.Run("simple lock", func(t *testing.T) {
		lock := New()
		mu := new(sync.Mutex)

		mu.Lock()
		i := atomic.NewInt64(0)
		go lock.Run("key #1", func() {
			i.Inc()
			mu.Lock()
		}, nil)

		go lock.Run("key #1", func() {
			i.Inc()
			mu.Lock()
		}, nil)

		time.Sleep(time.Millisecond * 100)
		assert.Equal(t, int64(1), i.Load())

		mu.Unlock()

		time.Sleep(time.Millisecond * 100)
		assert.Equal(t, int64(1), i.Load())
	})

	t.Run("high way to hell", func(t *testing.T) {
		wg := new(sync.WaitGroup)
		lock := New()

		count := 1000
		wg.Add(count)

		in := new(atomic.Int32)
		start := make(chan struct{})

		for i := 0; i < count; i++ {
			go func() {
				defer wg.Done()
				<-start
				lock.Run("test", func() {
					in.Inc()
					time.Sleep(time.Second)
				})
			}()
		}

		close(start)
		wg.Wait()

		assert.Equal(t, int32(1), in.Load())
	})
}
