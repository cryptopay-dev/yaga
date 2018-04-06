package workers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cryptopay-dev/yaga/logger/nop"
)

type mockLogger struct {
	*nop.Logger

	mu       *sync.Mutex
	messages []string
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		Logger: new(nop.Logger),
		mu:     new(sync.Mutex),
	}
}

func (l *mockLogger) Error(i ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, fmt.Sprint(i...))
}

func (l *mockLogger) Errorf(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, fmt.Sprintf(format, args...))
}

func (l *mockLogger) Contains(substr string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, msg := range l.messages {
		if strings.Contains(msg, substr) {
			return true
		}
	}
	return false
}
