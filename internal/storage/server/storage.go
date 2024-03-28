package storage

import "github.com/mrkovshik/yametrics/internal/model"

type (
	IServerStorage interface {
		UpdateMetricValue(newMetrics model.Metrics)
		GetMetricValue(newMetrics model.Metrics) (model.Metrics, error)
		GetAllMetrics() (string, error)
	}
)
