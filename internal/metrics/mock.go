package metrics

import (
	"context"
	"math/rand"
	"time"

	"github.com/mrkovshik/yametrics/internal/model"
)

// MockMetrics provides a mock implementation of MetricSource for testing purposes.
type MockMetrics struct {
	// MemStats represents mock memory statistics.
	MemStats map[string]float64
}

// NewMockMetrics creates a new instance of MockMetrics initialized with mock memory statistics.
func NewMockMetrics() MockMetrics {
	return MockMetrics{
		MemStats: map[string]float64{
			"Alloc":         1.00,
			"BuckHashSys":   2.00,
			"Frees":         3.00,
			"GCCPUFraction": 4.00,
		},
	}
}

// PollMemStats polls mock memory statistics and updates the provided storage.
// It implements the MetricSource interface.
func (m MockMetrics) PollMemStats(s storage) error {
	ctx := context.Background()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Update mock memory statistics as gauge metrics
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

	// Generate a random value and update it as a gauge metric
	randomValue := random.Float64()
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "RandomValue",
		MType: model.MetricTypeGauge,
		Value: &randomValue,
	}); err != nil {
		return err
	}

	// Increment poll count and update it as a counter metric
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
