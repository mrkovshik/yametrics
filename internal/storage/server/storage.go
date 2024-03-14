package storage

type (
	IServerStorage interface {
		UpdateMetricValue(metricType, metricName, metricValue string) error
		GetMetricValue(metricType, metricName string) (string, error)
		GetAllMetrics() (string, error)
	}
)
