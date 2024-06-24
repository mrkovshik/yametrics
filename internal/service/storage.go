package service

import (
	"context"

	"github.com/mrkovshik/yametrics/internal/model"
)

// Storage defines the interface for managing metrics data storage operations.
type Storage interface {
	// UpdateMetricValue updates or inserts a metric value in the storage.
	UpdateMetricValue(ctx context.Context, newMetrics model.Metrics) error

	// UpdateMetrics updates or inserts multiple metric values in the storage.
	UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error

	// GetMetricByModel retrieves a metric from the storage based on the provided model.
	GetMetricByModel(ctx context.Context, newMetrics model.Metrics) (model.Metrics, error)
	GetAllMetrics(ctx context.Context) (map[string]model.Metrics, error)
	StoreMetrics(ctx context.Context, path string) error

	// RestoreMetrics restores metrics from the specified path into the storage.
	RestoreMetrics(ctx context.Context, path string) error
	Ping(ctx context.Context) error
}
