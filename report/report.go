package report

import (
	"time"

	"github.com/cryptopay-dev/go-metrics"
	"github.com/labstack/echo"
)

// CaptureResponseTime used to track response timing
func CaptureResponseTime(platform, action string, now time.Time) {
	metrics.SendWithTags(metrics.M{
		"duration_ns": time.Since(now).Seconds() * 1e3,
	}, metrics.T{
		"platform": platform,
		"action":   action,
	}, "exchange_timings")
}

// CaptureResponseTimings used to track response code, endpoint and timing
func CaptureResponseTimings(res *echo.Response, endpoint string, now time.Time) {
	var code int

	if res != nil {
		code = res.Status
	}

	metrics.SendWithTags(metrics.M{
		"response_ns":   time.Since(now).Seconds() * 1e3,
		"response_code": code,
	}, metrics.T{
		"endpoint": endpoint,
	}, "response")
}

// CaptureMargin used to track margin for pair
func CaptureMargin(platform, pair string, margin float64) {
	metrics.SendWithTags(metrics.M{
		"margin": margin,
	}, metrics.T{
		"platform": platform,
		"pair":     pair,
	}, "margin")
}
