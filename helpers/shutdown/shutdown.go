package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type logger interface {
	Infof(string, ...interface{})
}

// ShutdownContext returns child context from passed context which will be canceled
// on incoming signals: SIGINT, SIGTERM, SIGHUP
func ShutdownContext(c context.Context, log logger) context.Context {
	ctx, cancel := context.WithCancel(c)
	go func() {
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		sig := <-ch
		if log != nil {
			log.Infof("received signal: %s", sig.String())
		}
	}()

	return ctx
}
