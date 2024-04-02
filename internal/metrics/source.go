package metrics

import (
	storage "github.com/mrkovshik/yametrics/internal/storage"
)

type MetricSource interface {
	PollMetrics(s storage.IStorage)
}
