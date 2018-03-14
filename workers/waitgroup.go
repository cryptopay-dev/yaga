package workers

import (
	"runtime"

	"go.uber.org/atomic"
)

type WaitGroup interface {
	Add(int)
	Done()
	Wait()
}

type waitGroup struct {
	count *atomic.Int64
}

func (wg waitGroup) Add(delta int) {
	if wg.count.Add(int64(delta)) < 0 {
		panic("waitgroup: negative WaitGroup counter")
	}
}

func (wg waitGroup) Done() {
	wg.Add(-1)
}

func (wg waitGroup) Wait() {
	for wg.count.Load() != 0 {
		runtime.Gosched()
	}
}

func NewWaitGroup() WaitGroup {
	return &waitGroup{
		count: atomic.NewInt64(0),
	}
}
