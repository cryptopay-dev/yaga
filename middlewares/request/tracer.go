package request

import (
	"github.com/cryptopay-dev/yaga/helpers"
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/cryptopay-dev/yaga/web"
)

const (
	RayTraceHeader = "X-Ray-Trace-ID"
)

type T = map[string]string

func rayTrace(ctx web.Context) (key, val string) {
	key = RayTraceHeader
	val = ctx.Request().Header.Get(key)
	return
}

func TraceTag(ctx web.Context) T {
	key, val := rayTrace(ctx)
	if val == "" {
		return nil
	}

	return T{key: val}
}

func RayTraceID(logger logger.Logger) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx web.Context) error {
			var (
				req = ctx.Request()
				res = ctx.Response()
				id  = req.Header.Get(RayTraceHeader)
			)

			if err := helpers.ValidateUUIDv4(id); err != nil {
				id = helpers.NewUUIDv4()
				req.Header.Set(RayTraceHeader, id)
			}

			res.Header().Set(RayTraceHeader, id)

			key, val := rayTrace(ctx)
			ctx.Echo().Logger = logger.WithContext(map[string]interface{}{key: val})

			return next(ctx)
		}
	}
}
