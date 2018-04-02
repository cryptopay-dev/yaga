package web

import (
	"github.com/cryptopay-dev/yaga/helpers"
	"github.com/cryptopay-dev/yaga/logger/log"
)

const (
	// RayTraceHeader key for headers
	RayTraceHeader = "X-Ray-Trace-ID"
)

// T is a tag
type T = map[string]string

// fetch ray-trace value from Request Header
func rayTrace(ctx Context) (key, val string) {
	key = RayTraceHeader
	val = ctx.Request().Header.Get(key)
	return
}

// TraceTag from Context
func TraceTag(ctx Context) T {
	key, val := rayTrace(ctx)
	if val == "" {
		return nil
	}

	return T{key: val}
}

// RayTraceID middleware
func RayTraceID(next HandlerFunc) HandlerFunc {
	return func(ctx Context) error {
		var (
			req = ctx.Request()
			res = ctx.Response()
			id  = req.Header.Get(RayTraceHeader)
		)

		if err := helpers.ValidateUUID(id); err != nil {
			id = helpers.NewUUID()
			req.Header.Set(RayTraceHeader, id)
		}

		res.Header().Set(RayTraceHeader, id)

		key, val := rayTrace(ctx)
		ctx.Echo().Logger = log.WithContext(map[string]interface{}{key: val})

		return next(ctx)
	}
}
