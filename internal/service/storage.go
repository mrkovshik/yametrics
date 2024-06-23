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

	// GetAllMetrics retrieves all metrics stored in the system as a formatted string.
	GetAllMetrics(ctx context.Context) (string, error)

	// StoreMetrics stores the current state of metrics to the specified path.
	StoreMetrics(ctx context.Context, path string) error

	// RestoreMetrics restores metrics from the specified path into the storage.
	RestoreMetrics(ctx context.Context, path string) error
}
