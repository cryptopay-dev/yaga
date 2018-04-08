package atomic

import (
	"github.com/cryptopay-dev/yaga/workers/locker"
	"go.uber.org/atomic"
)

type Lock struct {
	lock *atomic.Int32
}

func New() locker.Locker {
	return &Lock{
		lock: atomic.NewInt32(0),
	}
}

func (l *Lock) Run(key string, handler func()) {
	if !l.lock.CAS(0, 1) {
		return
	}
	defer l.lock.Store(0)
	handler()
}
