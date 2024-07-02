// Package service provides interfaces for managing metrics data storage operations.
package service

import (
	"context"

	"github.com/mrkovshik/yametrics/internal/model"
)

// Storage defines the interface for managing metrics data storage operations.
type Storage interface {
	// UpdateMetricValue updates or inserts a single metric value in the storage.
	// Parameters:
	// - ctx: the context to control the update operation.
	// - newMetrics: the Metrics model containing the metric data to be updated or inserted.
	// Returns:
	// - an error if the update operation fails.
	UpdateMetricValue(ctx context.Context, newMetrics model.Metrics) error

	// UpdateMetrics updates or inserts multiple metric values in the storage.
	// Parameters:
	// - ctx: the context to control the update operation.
	// - newMetrics: a slice of Metrics models containing the metric data to be updated or inserted.
	// Returns:
	// - an error if the update operation fails.
	UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error

	// GetMetricByModel retrieves a metric from the storage based on the provided model.
	// Parameters:
	// - ctx: the context to control the retrieval operation.
	// - newMetrics: the Metrics model specifying the metric to be retrieved.
	// Returns:
	// - the retrieved Metrics model.
	// - an error if the retrieval operation fails.
	GetMetricByModel(ctx context.Context, newMetrics model.Metrics) (model.Metrics, error)

	// GetAllMetrics retrieves all metrics from the storage.
	// Parameters:
	// - ctx: the context to control the retrieval operation.
	// Returns:
	// - a map of metric names to Metrics models representing all stored metrics.
	// - an error if the retrieval operation fails.
	GetAllMetrics(ctx context.Context) (map[string]model.Metrics, error)

	// StoreMetrics stores the current metrics to the specified path.
	// Parameters:
	// - ctx: the context to control the store operation.
	// - path: the file path where metrics should be stored.
	// Returns:
	// - an error if the store operation fails.
	StoreMetrics(ctx context.Context, path string) error

	// RestoreMetrics restores metrics from the specified path into the storage.
	// Parameters:
	// - ctx: the context to control the restore operation.
	// - path: the file path from where metrics should be restored.
	// Returns:
	// - an error if the restore operation fails.
	RestoreMetrics(ctx context.Context, path string) error

	// Ping checks the availability of the storage.
	// Parameters:
	// - ctx: the context to control the ping operation.
	// Returns:
	// - an error if the storage is not available.
	Ping(ctx context.Context) error
}
