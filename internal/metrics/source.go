package metrics

import (
	"github.com/mrkovshik/yametrics/internal/service"
)

// MetricSource defines the interface for metrics data sources.
type MetricSource interface {
	// PollMemStats polls memory statistics and updates the storage.
	PollMemStats(s service.Storage) error

	// PollVirtMemStats polls virtual memory statistics and updates the storage.
	PollVirtMemStats(s service.Storage) error
}
