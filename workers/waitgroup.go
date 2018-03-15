package workers

import (
	"sync"
)

type WaitGroup interface {
	Add(int)
	Done()
	Wait()
}

type waitGroup struct {
	count int64
	mu    *sync.RWMutex
	cond  *sync.Cond
}

func (wg *waitGroup) Add(delta int) {
	wg.mu.Lock()
	wg.count += int64(delta)
	if wg.count == 0 {
		wg.cond.Broadcast()
	} else if wg.count < 0 {
		wg.mu.Unlock()
		panic("waitgroup: negative WaitGroup counter")
	}
	wg.mu.Unlock()
}

func (wg *waitGroup) Done() {
	wg.Add(-1)
}

func (wg *waitGroup) Wait() {
	wg.cond.L.Lock()
	if wg.count > 0 {
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}

func NewWaitGroup() WaitGroup {
	mu := new(sync.RWMutex)
	return &waitGroup{
		cond: sync.NewCond(mu.RLocker()),
		mu:   mu,
	}
}
