package metrics

import (
	"github.com/mrkovshik/yametrics/internal/service"
)

type MetricSource interface {
	PollMemStats(s service.Storage) error
	PollVirtMemStats(s service.Storage) error
}
