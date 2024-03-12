package server

type (
	IStorage interface {
		UpdateCounter(Counter) error
		UpdateGauge(Gauge) error
		GetMetricValue(string, string) (string, error)
		GetAllMetrics() string
	}
)
