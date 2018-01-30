package report

import (
	"time"

	"github.com/cryptopay-dev/go-metrics"
)

func CaptureResponseTime(platform, action string, now time.Time) {
	metrics.SendWithTags(metrics.M{
		"duration_ns": time.Since(now).Seconds() * 1e3,
	}, metrics.T{
		"platform": platform,
		"action":   action,
	}, "exchange_timings")
}

func CaptureMargin(platform, pair string, margin float64) {
	metrics.SendWithTags(metrics.M{
		"margin": margin,
	}, metrics.T{
		"platform": platform,
		"pair":     pair,
	}, "margin")
}