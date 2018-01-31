package workers

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

var (
	minTickForTest = time.Microsecond * 10
	uniqWorkerN    = atomic.NewInt32(0)
)

func getUniqueWorkerName() string {
	return fmt.Sprintf("worker %d", uniqWorkerN.Inc())
}

func TestWorkerConflictName(t *testing.T) {
	c, creater := newCronForTest()
	defer c.StopCron()

	name := getUniqueWorkerName()
	w, err := creater(name, minTickForTest, func() {})
	if !assert.NoError(t, err) || !assert.NotNil(t, w) {
		assert.FailNow(t, "Cannot create worker")
	}
	if !assert.Equal(t, name, w.options.Name) || !assert.Equal(t, 1, len(w.pool.workers)) {
		assert.FailNow(t, "Invalidate workers data, must be 1 worker")
	}

	// create new worker with already name
	w, err = creater(name, minTickForTest, func() {})
	if !assert.Error(t, err) || !assert.Nil(t, w) {
		assert.FailNow(t, "Created new worker with duplicate name")
	}

	// create new worker with unique name
	w, err = creater(name+" foobar", minTickForTest, func() {})
	if !assert.NoError(t, err) || !assert.NotNil(t, w) {
		assert.FailNow(t, "Cannot create worker with unique name")
	}

	if !assert.Equal(t, 2, len(w.pool.workers)) {
		assert.FailNow(t, "Invalidate workers data, must be 2 workers")
	}
}

func TestWorkerStartAndStop(t *testing.T) {
	c, creater := newCronForTest()
	defer c.StopCron()

	t.Run("worker should be start", func(t *testing.T) {
		start := atomic.NewInt32(0)

		_, err := creater(getUniqueWorkerName(), minTickForTest, func() {
			start.Inc()
		})
		if !assert.NoError(t, err, "Cannot create worker") {
			t.FailNow()
		}

		for {
			// try to test start worker
			if start.Load() > 0 {
				break
			}
			runtime.Gosched()
		}
	})

	t.Run("worker should be stop", func(t *testing.T) {
		info := atomic.NewInt32(0)

		_, err := creater(getUniqueWorkerName(), minTickForTest, func() {
			info.Inc()
		})
		if !assert.NoError(t, err, "Cannot create worker") {
			t.FailNow()
		}

		for {
			// try to test start worker
			if info.Load() > 0 {
				break
			}
			runtime.Gosched()
		}

		c.Stop()

		c.Wait()
		info.Store(312)
		time.Sleep(minTickForTest * 100)

		for {
			// try to test stop worker
			if info.Load() == 312 {
				break
			}
			runtime.Gosched()
		}
	})
}

func TestWorkersRestart(t *testing.T) {
	c, creater := newCronForTest()
	defer c.StopCron()

	t.Run("worker should be restart", func(t *testing.T) {
		var (
			info       = atomic.NewInt32(0)
			num  int32 = 321
		)

		_, err := creater(getUniqueWorkerName(), minTickForTest, func() {
			info.Store(num)
		})
		if !assert.NoError(t, err, "Cannot create worker") {
			t.FailNow()
		}

		for {
			// try to test start worker
			if info.Load() == 321 {
				break
			}
			runtime.Gosched()
		}
		c.Stop()

		c.Wait()
		info.Store(1122)
		time.Sleep(minTickForTest * 100)

		for {
			// try to test stop worker
			if info.Load() == 1122 {
				break
			}
			runtime.Gosched()
		}

		num = 246975
		c.Start()

		for {
			// try to test restart worker
			if info.Load() == num {
				break
			}
			runtime.Gosched()
		}
	})

	t.Run("workers should be restart", func(t *testing.T) {
		info := atomic.NewInt32(0)

		_, err := creater(getUniqueWorkerName(), minTickForTest, func() {
			info.CAS(0, 11)
			info.CAS(456, 789)
		})
		if !assert.NoError(t, err, "Cannot create worker") {
			t.FailNow()
		}

		_, err = creater(getUniqueWorkerName(), minTickForTest, func() {
			info.CAS(11, 22)
			info.CAS(123, 456)
		})
		if !assert.NoError(t, err, "Cannot create worker") {
			t.FailNow()
		}

		for {
			// try to test start worker
			if info.Load() == 22 {
				break
			}
			runtime.Gosched()
		}

		c.Stop()

		c.Wait()
		info.Store(123)
		time.Sleep(minTickForTest * 100)

		for {
			// try to test stop workers
			if info.Load() == 123 {
				break
			}
			runtime.Gosched()
		}

		c.Start()

		for {
			// try to test start workers
			if info.Load() == 789 {
				break
			}
			runtime.Gosched()
		}
	})
}

func TestWorkersWait(t *testing.T) {
	c, creater := newCronForTest()
	defer c.StopCron()

	t.Run("workers should be wait while one worker locked", func(t *testing.T) {
		var (
			err error
			mu  sync.Mutex

			watch = make(chan struct{})
			info  = atomic.NewInt32(0)
		)

		mu.Lock()
		for i := 0; i < 5; i++ {
			lockedFlag := false
			if i == 4 {
				// only one worker should be block
				lockedFlag = true
			}
			n := int32(i)
			_, err = creater(getUniqueWorkerName(), minTickForTest, func() {
				info.CAS(n, n+1)
				if lockedFlag {
					lockedFlag = false
					mu.Lock()
				}
			})
			if !assert.NoError(t, err, "Cannot create worker") {
				t.FailNow()
			}
		}

		for {
			// try to test start workers
			if info.Load() == 5 {
				break
			}
			runtime.Gosched()
		}

		c.Stop()

		go func() {
			c.Wait()
			close(watch)
		}()

		select {
		case <-time.After(minTickForTest * 100):
		case <-watch:
			assert.FailNow(t, "Fail waiting of workers")
		}

		mu.Unlock() // unblock one worker
		<-watch
	})
}

func TestWorkersStop(t *testing.T) {
	c, creater := newCronForTest()
	defer c.StopCron()

	t.Run("all workers should be closed", func(t *testing.T) {
		var (
			err error
			num int32 = 2

			info = atomic.NewInt32(num)
		)

		for i := 0; i < 5; i++ {
			n := num
			_, err = creater(getUniqueWorkerName(), minTickForTest, func() {
				info.CAS(n, n*2)
				info.CAS(123, 75)
			})
			if !assert.NoError(t, err, "Cannot create worker") {
				t.FailNow()
			}

			num = num * 2
		}

		for {
			// try to test start workers
			if info.Load() == num {
				break
			}
			runtime.Gosched()
		}

		c.Stop()

		c.Wait()
		info.Store(123)
		time.Sleep(minTickForTest * 100)

		for {
			// try to test stop workers
			if info.Load() == 123 {
				break
			}
			runtime.Gosched()
		}
	})
}
