package errors

import (
	"errors"
	"fmt"

	"github.com/cryptopay-dev/yaga/middlewares/request"
	"github.com/cryptopay-dev/yaga/tracer"
	"github.com/getsentry/raven-go"
	"github.com/labstack/echo"
)

func (c *Logic) Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			defer func() {
				if rVal := recover(); rVal != nil {
					errorMsg := fmt.Sprint(rVal)
					err := errors.New(errorMsg)

					packet := tracer.StackPacket(err)
					raven.Capture(packet, request.TraceTag(ctx))
					ctx.Error(err)
				}
			}()

			return next(ctx)
		}
	}
}
