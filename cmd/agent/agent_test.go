package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrkovshik/yametrics/internal/metrics"
	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/storage"
)

func Test_getMetrics(t *testing.T) {
	var (
		src  = metrics.NewMockMetrics()
		strg = storage.NewMapStorage()
		ctx  = context.Background()
	)
	tests := []struct {
		name string
	}{
		{"positive 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err12 := src.PollMemStats(strg)
			assert.NoError(t, err12)
			PollCount, err1 := strg.GetMetricByModel(ctx, model.Metrics{
				ID:    "PollCount",
				MType: model.MetricTypeCounter,
			})
			assert.NoError(t, err1)
			assert.Equal(t, int64(1), *PollCount.Delta)
			Alloc, err2 := strg.GetMetricByModel(ctx, model.Metrics{
				ID:    "Alloc",
				MType: model.MetricTypeGauge,
			})
			assert.NoError(t, err2)
			assert.Equal(t, 1.00, *Alloc.Value)
			BuckHashSys, err3 := strg.GetMetricByModel(ctx, model.Metrics{
				ID:    "BuckHashSys",
				MType: model.MetricTypeGauge,
			})
			assert.NoError(t, err3)
			assert.Equal(t, 2.00, *BuckHashSys.Value)
			err13 := src.PollMemStats(strg)
			assert.NoError(t, err13)
			PollCount2, err11 := strg.GetMetricByModel(ctx, model.Metrics{
				ID:    "PollCount",
				MType: model.MetricTypeCounter,
			})
			assert.NoError(t, err11)
			assert.Equal(t, int64(2), *PollCount2.Delta)
		})
	}
}

func BenchmarkPollMemStats(b *testing.B) {
	var (
		src  = metrics.NewMockMetrics()
		strg = storage.NewMapStorage()
		ctx  = context.Background()
	)
	b.Run("poll", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			src.PollMemStats(strg)
		}
	})
	b.Run("get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strg.GetMetricByModel(ctx, model.Metrics{
				ID:    "BuckHashSys",
				MType: model.MetricTypeGauge,
			})
		}
	})
}
