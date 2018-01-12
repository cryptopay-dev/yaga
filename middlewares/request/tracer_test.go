package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var fakeRayTraceID = uuid.NewV4().String()

func TestRayTraceID(t *testing.T) {
	e := echo.New()
	e.Logger = nop.New()
	e.HideBanner = true
	req := httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var (
		err    error
		rec    = httptest.NewRecorder()
		c      = e.NewContext(req, rec)
		handle = RayTraceID(nop.New())(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})
	)

	assert.Nil(t, TraceTag(c))

	assert.False(t, traceIDSkipper(""))
	assert.False(t, traceIDSkipper("empty"))
	assert.False(t, traceIDSkipper(uuid.NewV1().String()))

	err = handle(c)
	assert.NoError(t, err)

	req.Header.Set(RayTraceHeader, fakeRayTraceID)
	handle = RayTraceID(nop.New())(func(c echo.Context) error {
		var (
			tag   = TraceTag(c)
			field = TraceField(c)
		)
		assert.Equal(t, c.Request().Header.Get(RayTraceHeader), fakeRayTraceID)
		assert.Equal(t, tag, T{RayTraceHeader: fakeRayTraceID})
		assert.Equal(t, field.Key, RayTraceHeader)
		assert.Equal(t, field.String, fakeRayTraceID)
		return c.NoContent(http.StatusOK)
	})

	err = handle(c)
	assert.NoError(t, err)
	assert.Equal(t, rec.Header().Get(RayTraceHeader), fakeRayTraceID)
}
