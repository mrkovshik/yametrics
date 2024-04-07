package metrics

import (
	"context"
	"math/rand"
	"time"

	"github.com/mrkovshik/yametrics/internal/model"

	"github.com/mrkovshik/yametrics/internal/storage"
)

type MockMetrics struct {
	MemStats map[string]float64
}

func NewMockMetrics() MockMetrics {
	return MockMetrics{
		map[string]float64{
			"Alloc":         1.00,
			"BuckHashSys":   2.00,
			"Frees":         3.00,
			"GCCPUFraction": 4.00,
		},
	}
}

func (m MockMetrics) PollMetrics(s storage.Storage) error {
	ctx := context.Background()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	alloc := 1.00
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "Alloc",
		MType: model.MetricTypeGauge,
		Value: &alloc,
	}); err != nil {
		return err
	}
	buckHashSys := 2.00
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "BuckHashSys",
		MType: model.MetricTypeGauge,
		Value: &buckHashSys,
	}); err != nil {
		return err
	}
	frees := 3.00
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "Frees",
		MType: model.MetricTypeGauge,
		Value: &frees,
	}); err != nil {
		return err
	}
	gCCPUFraction := 4.00
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "GCCPUFraction",
		MType: model.MetricTypeGauge,
		Value: &gCCPUFraction,
	}); err != nil {
		return err
	}
	randomValue := random.Float64()
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "RandomValue",
		MType: model.MetricTypeGauge,
		Value: &randomValue,
	}); err != nil {
		return err
	}
	delta := int64(1)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "PollCount",
		MType: model.MetricTypeCounter,
		Delta: &delta,
	}); err != nil {
		return err
	}
	return nil
}
