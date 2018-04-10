package workers

import (
	"github.com/cryptopay-dev/yaga/logger/nop"
	"go.uber.org/atomic"
)

type mockLogger struct {
	*nop.Logger
	count *atomic.Int32
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		Logger: new(nop.Logger),
		count:  atomic.NewInt32(0),
	}
}

func (l *mockLogger) Error(i ...interface{}) {
	l.count.Inc()
}

func (l *mockLogger) Errorf(format string, args ...interface{}) {
	l.count.Inc()
}

func (l *mockLogger) Count() int {
	return int(l.count.Load())
}
