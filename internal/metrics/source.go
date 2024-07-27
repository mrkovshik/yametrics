package metrics

import (
	"context"

	"github.com/mrkovshik/yametrics/internal/model"
)

type storage interface {
	UpdateMetricValue(ctx context.Context, newMetrics model.Metrics) error
	UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error
}

// MetricSource defines the interface for metrics data sources.
type MetricSource interface {
	// PollMemStats polls memory statistics and updates the storage.
	PollMemStats(s storage) error

	// PollVirtMemStats polls virtual memory statistics and updates the storage.
	PollVirtMemStats(s storage) error
}
