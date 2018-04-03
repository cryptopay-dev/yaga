package workers

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

var (
	limitTimeForTest = time.Second * 5
	minTickForTest   = time.Microsecond * 10
)

func checkGtZero(cnt *atomic.Int32) bool {
	limit := time.Now().Add(limitTimeForTest)

	for {
		if cnt.Load() > 0 {
			return true
		}

		if time.Now().After(limit) {
			return false
		}

		runtime.Gosched()
	}
}

func checkEqual(cnt *atomic.Int32, expected int32) bool {
	limit := time.Now().Add(limitTimeForTest)

	for {
		if cnt.Load() == expected {
			return true
		}

		if time.Now().After(limit) {
			return false
		}

		runtime.Gosched()
	}
}

func TestWorkerStartAndStop(t *testing.T) {
	w := New()
	defer w.Stop()

	t.Run("worker should be start", func(t *testing.T) {
		start := atomic.NewInt32(0)

		w.AddJob(minTickForTest, func(context.Context) {
			start.Inc()
		})

		if !checkGtZero(start) {
			assert.FailNow(t, "Cannot start worker")
		}
	})

	t.Run("worker should be stop", func(t *testing.T) {
		info := atomic.NewInt32(0)

		w.AddJob(minTickForTest, func(context.Context) {
			info.Inc()
		})

		if !checkGtZero(info) {
			assert.FailNow(t, "Cannot start worker")
		}

		w.Stop()

		w.Wait(nil)
		info.Store(312)
		time.Sleep(minTickForTest * 100)

		if !checkEqual(info, 312) {
			assert.FailNow(t, "Cannot stop worker")
		}
	})
}

func TestWorkersContext(t *testing.T) {
	w := New()
	defer w.Stop()

	t.Run("workers should be stop when context done", func(t *testing.T) {
		info := atomic.NewInt32(0)

		for i := 0; i < 5; i++ {
			n := int32(i)
			w.AddJob(minTickForTest, func(ctx context.Context) {
				if info.CAS(n, n+1) {
					<-ctx.Done()
					info.Add(2)
				}
			})
		}

		if !checkEqual(info, 5) {
			assert.FailNow(t, "Cannot start workers")
		}

		ctx, cancel := context.WithTimeout(context.Background(), minTickForTest*100)
		defer cancel()
		err := w.Wait(ctx)
		if !assert.Error(t, err, "Fail waiting of workers") {
			t.FailNow()
		}

		w.Stop()

		if !checkEqual(info, 15) {
			assert.FailNow(t, "Context done failed")
		}

		w.Wait(nil)
	})
}

func TestWorkersWait(t *testing.T) {
	w := New()
	defer w.Stop()

	t.Run("workers should be wait while one worker locked", func(t *testing.T) {
		var (
			mu sync.Mutex

			info = atomic.NewInt32(0)
		)

		mu.Lock()
		for i := 0; i < 5; i++ {
			lockedFlag := atomic.NewBool(false)
			if i == 4 {
				// we will block only one worker
				lockedFlag.Store(true)
			}
			n := int32(i)
			w.AddJob(minTickForTest, func(context.Context) {
				info.CAS(n, n+1)
				if lockedFlag.Load() {
					lockedFlag.Store(false)
					mu.Lock()
				}
			})
		}

		if !checkEqual(info, 5) {
			assert.FailNow(t, "Cannot start workers")
		}

		w.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), minTickForTest*100)
		defer cancel()
		err := w.Wait(ctx)
		if !assert.Error(t, err) {
			assert.FailNow(t, "Fail waiting of workers")
		}

		mu.Unlock() // unblock one worker
		w.Wait(nil)
	})
}

// TODO need?
func TestWorkersWithError(t *testing.T) {
	w := New()
	defer w.Stop()

	t.Run("workers should be stop when one worker returns error", func(t *testing.T) {
		var mu sync.Mutex
		errWorker := errors.New("test worker error")
		info := atomic.NewInt32(0)

		mu.Lock()
		for i := 0; i < 5; i++ {
			n := int32(i)
			w.AddJobWithError(minTickForTest, func(ctx context.Context) error {
				if info.CAS(n, n+1) {
					if n == 2 {
						// we will block only one worker
						mu.Lock()
						return errWorker
					}
					<-ctx.Done()
					info.Add(2)
				}
				return nil
			})
		}

		if !checkEqual(info, 5) {
			assert.FailNow(t, "Cannot start workers")
		}

		ctx, cancel := context.WithTimeout(context.Background(), minTickForTest*100)
		defer cancel()
		err := w.Wait(ctx)
		if !assert.Error(t, err, "Fail waiting of workers") {
			t.FailNow()
		}

		mu.Unlock() // unblock one worker

		if !checkEqual(info, 13) {
			assert.FailNow(t, "Context done failed")
		}

		err = w.Wait(nil)
		if !assert.Error(t, err, "Failed to return worker error") {
			t.FailNow()
		}
		assert.Contains(t, err.Error(), errWorker.Error())
	})
}

func TestWorkersStop(t *testing.T) {
	w := New()
	defer w.Stop()

	t.Run("all workers should be closed", func(t *testing.T) {
		var (
			num int32 = 2

			info = atomic.NewInt32(num)
		)

		for i := 0; i < 5; i++ {
			n := num
			w.AddJob(minTickForTest, func(context.Context) {
				info.CAS(n, n*2)
				info.CAS(123, 75)
			})

			num = num * 2
		}

		if !checkEqual(info, num) {
			assert.FailNow(t, "Cannot start workers")
		}

		w.Stop()

		w.Wait(nil)
		info.Store(123)
		time.Sleep(minTickForTest * 100)

		if !checkEqual(info, 123) {
			assert.FailNow(t, "Cannot stop workers")
		}
	})
}
