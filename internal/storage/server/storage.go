package server

type (
	IStorage interface {
		UpdateCounter(Counter) error
		UpdateGauge(Gauge) error
		GetCounterValue(string) (string, error)
		GetGaugeValue(string) (string, error)
		GetAllMetrics() string
	}
)
