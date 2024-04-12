package metrics

import (
	"github.com/mrkovshik/yametrics/internal/service"
)

type MetricSource interface {
	PollMetrics(s service.Storage) error
}
