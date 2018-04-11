package web

import (
	"errors"
	"fmt"

	"github.com/cryptopay-dev/yaga/tracer"
	"github.com/getsentry/raven-go"
	"github.com/labstack/echo"
)

// Recover is an echo-middleware to capture panics in controllers/actions
// and send info to sentry
func (c *Logic) Recover() MiddlewareFunc {
	return func(next HandlerFunc) echo.HandlerFunc {
		return func(ctx Context) error {
			defer func() {
				if rVal := recover(); rVal != nil {
					errorMsg := fmt.Sprint(rVal)
					err := errors.New(errorMsg)

					packet := tracer.StackPacket(err)
					raven.Capture(packet, TraceTag(ctx))
					ctx.Error(err)
				}
			}()

			return next(ctx)
		}
	}
}
