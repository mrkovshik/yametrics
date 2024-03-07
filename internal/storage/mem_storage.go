package storage

type (
	MapStorage struct {
		Gauges   map[string]float64
		Counters map[string]int64
	}
	IStorage interface {
		UpdateCounter(Counter) error
		UpdateGauge(Gauge) error
		GetCounterValue(string) string
		GetGaugeValue(string) string
		GetAllMetrics() string
	}
)
