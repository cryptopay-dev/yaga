package metrics

// Provider interface describe all accessible variants for to write metrics
type Provider interface {
	// IncrementCounter increase counter metric,
	// a counter is a cumulative metric that represents a single numerical value that only ever goes up
	IncrementCounter(key string, value uint)
	// IncrementGauge add value to gauge metric,
	// a gauge is a metric that represents a single numerical value that can arbitrarily go up and down
	IncrementGauge(key string, value int)
	// Observe write new value to summaray metric
	Observe(key string, value float64)
}

var defaultProvider Provider = &nop{}

// SetProvider set global metrics provider
func SetProvider(p Provider) {
	defaultProvider = p
}

// IncrementCounter increase key counter
func IncrementCounter(key string, value uint) {
	defaultProvider.IncrementCounter(key, value)
}

// IncrementGauge add value to gauge by key
func IncrementGauge(key string, value int) {
	defaultProvider.IncrementGauge(key, value)
}

// Observe write new value to summaray metric
func Observe(key string, value float64) {
	defaultProvider.Observe(key, value)
}
