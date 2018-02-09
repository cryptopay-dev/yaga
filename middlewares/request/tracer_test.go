package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptopay-dev/yaga/helpers"
	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/web"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var fakeRayTraceID = helpers.GenerateUUIDv4AsString()

func TestRayTraceID(t *testing.T) {
	e := web.New(web.Options{})

	req := httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var (
		err    error
		rec    = httptest.NewRecorder()
		c      = e.NewContext(req, rec)
		handle = RayTraceID(nop.New())(func(c web.Context) error {
			return c.NoContent(http.StatusOK)
		})
	)

	assert.Nil(t, TraceTag(c))

	err = handle(c)
	assert.NoError(t, err)

	req.Header.Set(RayTraceHeader, fakeRayTraceID)
	handle = RayTraceID(nop.New())(func(c echo.Context) error {
		var (
			tag   = TraceTag(c)
			field = TraceTag(c)
		)
		assert.Equal(t, c.Request().Header.Get(RayTraceHeader), fakeRayTraceID)
		assert.Equal(t, tag, T{RayTraceHeader: fakeRayTraceID})
		assert.Equal(t, field["X-Ray-Trace-ID"], fakeRayTraceID)
		return c.NoContent(http.StatusOK)
	})

	err = handle(c)
	assert.NoError(t, err)
	assert.Equal(t, rec.Header().Get(RayTraceHeader), fakeRayTraceID)
}
