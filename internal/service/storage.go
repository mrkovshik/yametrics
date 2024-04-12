package service

import (
	"context"

	"github.com/mrkovshik/yametrics/internal/model"
)

type Storage interface {
	UpdateMetricValue(ctx context.Context, newMetrics model.Metrics) error
	UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error
	GetMetricByModel(ctx context.Context, newMetrics model.Metrics) (model.Metrics, error)
	GetAllMetrics(ctx context.Context) (string, error)
	StoreMetrics(ctx context.Context, path string) error
	RestoreMetrics(ctx context.Context, path string) error
}
