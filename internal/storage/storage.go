package storage

import "github.com/mrkovshik/yametrics/internal/model"

type (
	Storage interface {
		UpdateMetricValue(newMetrics model.Metrics)
		GetMetricByModel(newMetrics model.Metrics) (model.Metrics, error)
		GetAllMetrics() (string, error)
		StoreMetrics(path string) error
		RestoreMetrics(path string) error
	}
)
