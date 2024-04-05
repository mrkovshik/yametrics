package metrics

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/mrkovshik/yametrics/internal/model"

	"github.com/mrkovshik/yametrics/internal/storage"
)

type RuntimeMetrics struct {
	MemStats runtime.MemStats
}

func NewRuntimeMetrics() RuntimeMetrics {
	m := RuntimeMetrics{
		MemStats: runtime.MemStats{},
	}
	return m
}

func (m RuntimeMetrics) PollMetrics(s storage.Storage) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	runtime.ReadMemStats(&m.MemStats)
	valueAlloc := float64(m.MemStats.Alloc)
	s.UpdateMetricValue(model.Metrics{
		ID:    "Alloc",
		MType: model.MetricTypeGauge,
		Value: &valueAlloc,
	})
	valueBuckHashSys := float64(m.MemStats.BuckHashSys)
	s.UpdateMetricValue(model.Metrics{
		ID:    "BuckHashSys",
		MType: model.MetricTypeGauge,
		Value: &valueBuckHashSys,
	})
	valueFrees := float64(m.MemStats.Frees)
	s.UpdateMetricValue(model.Metrics{
		ID:    "Frees",
		MType: model.MetricTypeGauge,
		Value: &valueFrees,
	})
	s.UpdateMetricValue(model.Metrics{
		ID:    "GCCPUFraction",
		MType: model.MetricTypeGauge,
		Value: &m.MemStats.GCCPUFraction,
	})
	valueGCSys := float64(m.MemStats.GCSys)
	s.UpdateMetricValue(model.Metrics{
		ID:    "GCSys",
		MType: model.MetricTypeGauge,
		Value: &valueGCSys,
	})
	valueHeapAlloc := float64(m.MemStats.HeapAlloc)
	s.UpdateMetricValue(model.Metrics{
		ID:    "HeapAlloc",
		MType: model.MetricTypeGauge,
		Value: &valueHeapAlloc,
	})
	valueHeapIdle := float64(m.MemStats.HeapIdle)
	s.UpdateMetricValue(model.Metrics{
		ID:    "HeapIdle",
		MType: model.MetricTypeGauge,
		Value: &valueHeapIdle,
	})
	valueHeapInuse := float64(m.MemStats.HeapInuse)
	s.UpdateMetricValue(model.Metrics{
		ID:    "HeapInuse",
		MType: model.MetricTypeGauge,
		Value: &valueHeapInuse,
	})
	valueHeapObjects := float64(m.MemStats.HeapObjects)
	s.UpdateMetricValue(model.Metrics{
		ID:    "HeapObjects",
		MType: model.MetricTypeGauge,
		Value: &valueHeapObjects,
	})
	valueHeapReleased := float64(m.MemStats.HeapReleased)
	s.UpdateMetricValue(model.Metrics{
		ID:    "HeapReleased",
		MType: model.MetricTypeGauge,
		Value: &valueHeapReleased,
	})
	valueHeapSys := float64(m.MemStats.HeapSys)
	s.UpdateMetricValue(model.Metrics{
		ID:    "HeapSys",
		MType: model.MetricTypeGauge,
		Value: &valueHeapSys,
	})
	valueLastGC := float64(m.MemStats.LastGC)
	s.UpdateMetricValue(model.Metrics{
		ID:    "LastGC",
		MType: model.MetricTypeGauge,
		Value: &valueLastGC,
	})
	valueLookups := float64(m.MemStats.Lookups)
	s.UpdateMetricValue(model.Metrics{
		ID:    "Lookups",
		MType: model.MetricTypeGauge,
		Value: &valueLookups,
	})
	valueMCacheInuse := float64(m.MemStats.MCacheInuse)
	s.UpdateMetricValue(model.Metrics{
		ID:    "MCacheInuse",
		MType: model.MetricTypeGauge,
		Value: &valueMCacheInuse,
	})
	valueMCacheSys := float64(m.MemStats.MCacheSys)
	s.UpdateMetricValue(model.Metrics{
		ID:    "MCacheSys",
		MType: model.MetricTypeGauge,
		Value: &valueMCacheSys,
	})
	valueMSpanInuse := float64(m.MemStats.MSpanInuse)
	s.UpdateMetricValue(model.Metrics{
		ID:    "MSpanInuse",
		MType: model.MetricTypeGauge,
		Value: &valueMSpanInuse,
	})
	valueMSpanSys := float64(m.MemStats.MSpanSys)
	s.UpdateMetricValue(model.Metrics{
		ID:    "MSpanSys",
		MType: model.MetricTypeGauge,
		Value: &valueMSpanSys,
	})
	valueMallocs := float64(m.MemStats.Mallocs)
	s.UpdateMetricValue(model.Metrics{
		ID:    "Mallocs",
		MType: model.MetricTypeGauge,
		Value: &valueMallocs,
	})
	valueNextGC := float64(m.MemStats.NextGC)
	s.UpdateMetricValue(model.Metrics{
		ID:    "NextGC",
		MType: model.MetricTypeGauge,
		Value: &valueNextGC,
	})
	valueNumForcedGC := float64(m.MemStats.NumForcedGC)
	s.UpdateMetricValue(model.Metrics{
		ID:    "NumForcedGC",
		MType: model.MetricTypeGauge,
		Value: &valueNumForcedGC,
	})
	valueNumGC := float64(m.MemStats.NumGC)
	s.UpdateMetricValue(model.Metrics{
		ID:    "NumGC",
		MType: model.MetricTypeGauge,
		Value: &valueNumGC,
	})
	valueOtherSys := float64(m.MemStats.OtherSys)
	s.UpdateMetricValue(model.Metrics{
		ID:    "OtherSys",
		MType: model.MetricTypeGauge,
		Value: &valueOtherSys,
	})
	valuePauseTotalNs := float64(m.MemStats.PauseTotalNs)
	s.UpdateMetricValue(model.Metrics{
		ID:    "PauseTotalNs",
		MType: model.MetricTypeGauge,
		Value: &valuePauseTotalNs,
	})
	valueStackInuse := float64(m.MemStats.StackInuse)
	s.UpdateMetricValue(model.Metrics{
		ID:    "StackInuse",
		MType: model.MetricTypeGauge,
		Value: &valueStackInuse,
	})
	valueStackSys := float64(m.MemStats.StackSys)
	s.UpdateMetricValue(model.Metrics{
		ID:    "StackSys",
		MType: model.MetricTypeGauge,
		Value: &valueStackSys,
	})
	valueSys := float64(m.MemStats.Sys)
	s.UpdateMetricValue(model.Metrics{
		ID:    "Sys",
		MType: model.MetricTypeGauge,
		Value: &valueSys,
	})
	valueTotalAlloc := float64(m.MemStats.TotalAlloc)
	s.UpdateMetricValue(model.Metrics{
		ID:    "TotalAlloc",
		MType: model.MetricTypeGauge,
		Value: &valueTotalAlloc,
	})
	valueRandom := random.Float64()
	s.UpdateMetricValue(model.Metrics{
		ID:    "RandomValue",
		MType: model.MetricTypeGauge,
		Value: &valueRandom,
	})
	delta := int64(1)
	s.UpdateMetricValue(model.Metrics{
		ID:    "PollCount",
		MType: model.MetricTypeCounter,
		Delta: &delta,
	})
}
