package api

import (
	"context"

	"github.com/mrkovshik/yametrics/internal/model"
)

// Service represents an interface for managing metrics.
type Service interface {
	// UpdateMetrics updates a batch of metrics.
	// Parameters:
	// - ctx: the context to control the update operation.
	// - batch: a slice of Metrics to be updated.
	// Returns:
	// - an error if the update operation fails.
	UpdateMetrics(ctx context.Context, batch []model.Metrics) error

	// GetMetric retrieves a single metric based on the provided metric model.
	// Parameters:
	// - ctx: the context to control the retrieval operation.
	// - metricModel: the model representing the metric to be retrieved.
	// Returns:
	// - the retrieved Metrics.
	// - an error if the retrieval operation fails.
	GetMetric(ctx context.Context, metricModel model.Metrics) (model.Metrics, error)

	// GetAllMetrics retrieves all available metrics.
	// Parameters:
	// - ctx: the context to control the retrieval operation.
	// Returns:
	// - a string representing all metrics.
	// - an error if the retrieval operation fails.
	GetAllMetrics(ctx context.Context) (string, error)

	// Ping checks the availability of the service.
	// Parameters:
	// - ctx: the context to control the ping operation.
	// Returns:
	// - an error if the service is not available.
	Ping(ctx context.Context) error
}
