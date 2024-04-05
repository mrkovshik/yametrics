package metrics

import (
	"github.com/mrkovshik/yametrics/internal/storage"
)

type MetricSource interface {
	PollMetrics(s storage.Storage)
}
