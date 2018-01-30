package request

import (
	"github.com/cryptopay-dev/yaga/logger"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

const (
	RayTraceHeader = "X-Ray-Trace-ID"
)

type T = map[string]string

func rayTrace(ctx echo.Context) (key, val string) {
	key = RayTraceHeader
	val = ctx.Request().Header.Get(key)
	return
}

func TraceTag(ctx echo.Context) T {
	key, val := rayTrace(ctx)
	if val == "" {
		return nil
	}

	return T{key: val}
}

func RayTraceID(logger logger.Logger) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var (
				req = ctx.Request()
				res = ctx.Response()
				id  = req.Header.Get(RayTraceHeader)
			)

			if !traceIDSkipper(id) {
				id = traceIDGenerator()
				req.Header.Set(RayTraceHeader, id)
			}

			res.Header().Set(RayTraceHeader, id)

			key, val := rayTrace(ctx)
			ctx.Echo().Logger = logger.WithContext(map[string]interface{}{key: val})

			return next(ctx)
		}
	}
}

func traceIDSkipper(id string) bool {
	if id == "" {
		return false
	} else if uid, err := uuid.FromString(id); err != nil {
		return false
	} else if uid.Version() != 4 {
		return false
	}

	return true
}

func traceIDGenerator() string {
	return uuid.NewV4().String()
}
