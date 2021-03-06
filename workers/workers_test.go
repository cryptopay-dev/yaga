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
	limitTimeForTest = time.Second * 5
	minTickForTest   = time.Microsecond * 10
	uniqWorkerN      = atomic.NewInt32(0)
)

func getUniqueWorkerName() string {
	return fmt.Sprintf("worker %d", uniqWorkerN.Inc())
}

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

	// create new worker with existing name
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
	c.Start()

	t.Run("worker should be start", func(t *testing.T) {
		start := atomic.NewInt32(0)

		_, err := creater(getUniqueWorkerName(), minTickForTest, func() {
			start.Inc()
		})
		if !assert.NoError(t, err, "Cannot create worker") {
			t.FailNow()
		}

		if !checkGtZero(start) {
			assert.FailNow(t, "Cannot start worker")
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

		if !checkGtZero(info) {
			assert.FailNow(t, "Cannot start worker")
		}

		c.Stop()

		c.Wait()
		info.Store(312)
		time.Sleep(minTickForTest * 100)

		if !checkEqual(info, 312) {
			assert.FailNow(t, "Cannot stop worker")
		}
	})
}

func TestWorkersRestart(t *testing.T) {
	c, creater := newCronForTest()
	defer c.StopCron()
	c.Start()

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

		if !checkEqual(info, 321) {
			assert.FailNow(t, "Cannot start worker")
		}
		c.Stop()

		c.Wait()
		info.Store(1122)
		time.Sleep(minTickForTest * 100)

		if !checkEqual(info, 1122) {
			assert.FailNow(t, "Cannot stop worker")
		}

		num = 246975
		c.Start()

		if !checkEqual(info, num) {
			assert.FailNow(t, "Cannot restart worker")
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

		if !checkEqual(info, 22) {
			assert.FailNow(t, "Cannot start workers")
		}

		c.Stop()

		c.Wait()
		info.Store(123)
		time.Sleep(minTickForTest * 100)

		if !checkEqual(info, 123) {
			assert.FailNow(t, "Cannot stop workers")
		}

		c.Start()

		if !checkEqual(info, 789) {
			assert.FailNow(t, "Cannot start workers")
		}
	})
}

func TestWorkersWait(t *testing.T) {
	c, creater := newCronForTest()
	defer c.StopCron()
	c.Start()

	t.Run("workers should be wait while one worker locked", func(t *testing.T) {
		var (
			err error
			mu  sync.Mutex

			watch = make(chan struct{})
			info  = atomic.NewInt32(0)
		)

		mu.Lock()
		for i := 0; i < 5; i++ {
			lockedFlag := atomic.NewBool(false)
			if i == 4 {
				// we will block only one worker
				lockedFlag.Store(true)
			}
			n := int32(i)
			_, err = creater(getUniqueWorkerName(), minTickForTest, func() {
				info.CAS(n, n+1)
				if lockedFlag.Load() {
					lockedFlag.Store(false)
					mu.Lock()
				}
			})
			if !assert.NoError(t, err, "Cannot create worker") {
				t.FailNow()
			}
		}

		if !checkEqual(info, 5) {
			assert.FailNow(t, "Cannot start workers")
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
	c.Start()

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

		if !checkEqual(info, num) {
			assert.FailNow(t, "Cannot start workers")
		}

		c.Stop()

		c.Wait()
		info.Store(123)
		time.Sleep(minTickForTest * 100)

		if !checkEqual(info, 123) {
			assert.FailNow(t, "Cannot stop workers")
		}
	})
}
