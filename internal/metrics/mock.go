package metrics

import (
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

func (m MockMetrics) PollMetrics(s storage.Storage) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	alloc := 1.00
	s.UpdateMetricValue(model.Metrics{
		ID:    "Alloc",
		MType: model.MetricTypeGauge,
		Value: &alloc,
	})
	buckHashSys := 2.00
	s.UpdateMetricValue(model.Metrics{
		ID:    "BuckHashSys",
		MType: model.MetricTypeGauge,
		Value: &buckHashSys,
	})
	frees := 3.00
	s.UpdateMetricValue(model.Metrics{
		ID:    "Frees",
		MType: model.MetricTypeGauge,
		Value: &frees,
	})
	gCCPUFraction := 4.00
	s.UpdateMetricValue(model.Metrics{
		ID:    "GCCPUFraction",
		MType: model.MetricTypeGauge,
		Value: &gCCPUFraction,
	})
	randomValue := random.Float64()
	s.UpdateMetricValue(model.Metrics{
		ID:    "RandomValue",
		MType: model.MetricTypeGauge,
		Value: &randomValue,
	})
	delta := int64(1)
	s.UpdateMetricValue(model.Metrics{
		ID:    "PollCount",
		MType: model.MetricTypeCounter,
		Delta: &delta,
	})
}
