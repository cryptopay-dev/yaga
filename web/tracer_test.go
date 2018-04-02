package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptopay-dev/yaga/helpers"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var fakeRayTraceID = helpers.NewUUID()

func TestRayTraceID(t *testing.T) {
	e, err := New(Options{})
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var (
		rec    = httptest.NewRecorder()
		c      = e.NewContext(req, rec)
		handle = RayTraceID(func(c Context) error {
			return c.NoContent(http.StatusOK)
		})
	)

	assert.Nil(t, TraceTag(c))

	err = handle(c)
	assert.NoError(t, err)

	req.Header.Set(RayTraceHeader, fakeRayTraceID)
	handle = RayTraceID(func(c echo.Context) error {
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
