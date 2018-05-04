package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()
}

// Context signal-bound context
func Context() context.Context {
	return ctx
}

// Cancel context
func Cancel() {
	cancel()
}

// Wait for context done
func Wait() error {
	<-ctx.Done()
	return ctx.Err()
}
