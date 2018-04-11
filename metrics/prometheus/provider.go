package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Provider is addapter from prometheus client to metrics.Provider interface
type Provider struct {
	server *http.Server

	Counter *prometheus.CounterVec
	Gauge   *prometheus.GaugeVec
	Summary *prometheus.SummaryVec
}

// NewProvider yield new prometheus providert instance
func NewProvider(bindAddress string) *Provider {
	p := &Provider{
		server: &http.Server{
			Addr:    bindAddress,
			Handler: promhttp.Handler(),
		},

		Counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "app_counters",
				Help: "is common counter metrics with different `key` label",
			},
			[]string{"key"},
		),
		Gauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "app_gauages",
				Help: "is common gauage metrics with different `key` label",
			},
			[]string{"key"},
		),
		Summary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "app_summaries",
				Help:       "is common summary metrics with different `key` label",
				Objectives: map[float64]float64{0.5: 0.5, 0.75: 0.75, 0.95: 0.95, 0.99: 0.99},
			},
			[]string{"key"},
		),
	}

	prometheus.Register(p.Counter)
	prometheus.Register(p.Gauge)
	prometheus.Register(p.Summary)

	return p
}

// IncrementCounter increase key counter
func (p Provider) IncrementCounter(key string, value uint) {
	p.Counter.WithLabelValues(key).Add(float64(value))
}

// IncrementGauge add value to gauge by key
func (p Provider) IncrementGauge(key string, value int) {
	p.Gauge.WithLabelValues(key).Add(float64(value))
}

// Observe write new value to summaray metric
func (p Provider) Observe(key string, value float64) {
	p.Summary.WithLabelValues(key).Observe(value)
}
