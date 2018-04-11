package local

import (
	"sync"

	"github.com/cryptopay-dev/yaga/locker"
)

type Lock struct {
	mu   *sync.Mutex
	keys map[string]bool
}

func New() locker.Locker {
	return &Lock{
		mu:   new(sync.Mutex),
		keys: make(map[string]bool),
	}
}

func (l *Lock) unlock(key string) {
	l.mu.Lock()
	l.keys[key] = false
	l.mu.Unlock()
}

func (l *Lock) Run(key string, handler func(), options ...locker.Option) error {
	l.mu.Lock()
	lock := l.keys[key]
	if lock {
		l.mu.Unlock()
		return nil
	}
	l.keys[key] = true
	l.mu.Unlock()

	defer l.unlock(key)
	handler()

	return nil
}
