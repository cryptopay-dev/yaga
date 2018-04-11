package metrics

type nop struct{}

func (n nop) IncrementCounter(key string, value uint) {}
func (n nop) IncrementGauge(key string, value int)    {}
func (n nop) Observe(key string, value float64)       {}
