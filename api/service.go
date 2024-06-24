package api

import (
	"context"

	"github.com/mrkovshik/yametrics/internal/model"
)

type Service interface {
	UpdateMetrics(ctx context.Context, batch []model.Metrics) error
	GetMetric(ctx context.Context, metricModel model.Metrics) (model.Metrics, error)
	GetAllMetrics(ctx context.Context) (string, error)
	Ping(ctx context.Context) error
}
