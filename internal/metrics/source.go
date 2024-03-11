package metrics

import (
	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
)

type MetricSource interface {
	PollMetrics(s storage.IAgentStorage)
}
