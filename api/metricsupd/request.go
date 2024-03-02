package metricsupd

type UpdateMetricRequest struct {
	metricType  string
	metricName  string
	metricValue int64
}
