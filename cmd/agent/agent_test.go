package main

import (
	"github.com/mrkovshik/yametrics/internal/model"
	"testing"

	"github.com/mrkovshik/yametrics/internal/metrics"
	storage "github.com/mrkovshik/yametrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func Test_getMetrics(t *testing.T) {
	var (
		src  = metrics.NewMockMetrics()
		strg = storage.NewMapStorage()
	)
	tests := []struct {
		name string
	}{
		{"positive 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src.PollMetrics(strg)
			PollCount, err1 := strg.GetMetricByModel(model.Metrics{
				ID:    "PollCount",
				MType: model.MetricTypeCounter,
			})
			assert.NoError(t, err1)
			assert.Equal(t, int64(1), *PollCount.Delta)
			Alloc, err2 := strg.GetMetricByModel(model.Metrics{
				ID:    "Alloc",
				MType: model.MetricTypeGauge,
			})
			assert.NoError(t, err2)
			assert.Equal(t, 1.00, *Alloc.Value)
			BuckHashSys, err3 := strg.GetMetricByModel(model.Metrics{
				ID:    "BuckHashSys",
				MType: model.MetricTypeGauge,
			})
			assert.NoError(t, err3)
			assert.Equal(t, 2.00, *BuckHashSys.Value)
			src.PollMetrics(strg)
			PollCount2, err11 := strg.GetMetricByModel(model.Metrics{
				ID:    "PollCount",
				MType: model.MetricTypeCounter,
			})
			assert.NoError(t, err11)
			assert.Equal(t, int64(2), *PollCount2.Delta)
		})
	}
}
