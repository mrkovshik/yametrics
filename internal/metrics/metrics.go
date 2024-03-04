package metrics

import (
	"github.com/mrkovshik/yametrics/internal/storage"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type (
	Imetric interface {
		Update(storage.IStorage) error
	}
	MetricSource interface {
		GetMetrics(MetricsValues *sync.Map)
	}
	RuntimeMetrics struct {
		MemStats runtime.MemStats
	}
	MockMetrics struct {
		MemStats map[string]float64
	}
)

func NewRuntimeMetrics() RuntimeMetrics {
	m := RuntimeMetrics{
		MemStats: runtime.MemStats{},
	}
	runtime.ReadMemStats(&m.MemStats)
	return m
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

func (m RuntimeMetrics) GetMetrics(MetricsValues *sync.Map) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	MetricsValues.Store("Alloc", float64(m.MemStats.Alloc))
	MetricsValues.Store("BuckHashSys", float64(m.MemStats.BuckHashSys))
	MetricsValues.Store("Frees", float64(m.MemStats.Frees))
	MetricsValues.Store("GCCPUFraction", m.MemStats.GCCPUFraction)
	MetricsValues.Store("GCSys", float64(m.MemStats.GCSys))
	MetricsValues.Store("HeapAlloc", float64(m.MemStats.HeapAlloc))
	MetricsValues.Store("HeapIdle", float64(m.MemStats.HeapIdle))
	MetricsValues.Store("HeapInuse", float64(m.MemStats.HeapInuse))
	MetricsValues.Store("HeapObjects", float64(m.MemStats.HeapObjects))
	MetricsValues.Store("HeapReleased", float64(m.MemStats.HeapReleased))
	MetricsValues.Store("HeapSys", float64(m.MemStats.HeapSys))
	MetricsValues.Store("LastGC", float64(m.MemStats.LastGC))
	MetricsValues.Store("Lookups", float64(m.MemStats.Lookups))
	MetricsValues.Store("MCacheInuse", float64(m.MemStats.MCacheInuse))
	MetricsValues.Store("MCacheSys", float64(m.MemStats.MCacheSys))
	MetricsValues.Store("MSpanInuse", float64(m.MemStats.MSpanInuse))
	MetricsValues.Store("MSpanSys", float64(m.MemStats.MSpanSys))
	MetricsValues.Store("Mallocs", float64(m.MemStats.Mallocs))
	MetricsValues.Store("NextGC", float64(m.MemStats.NextGC))
	MetricsValues.Store("NumForcedGC", float64(m.MemStats.NumForcedGC))
	MetricsValues.Store("NumGC", float64(m.MemStats.NumGC))
	MetricsValues.Store("OtherSys", float64(m.MemStats.OtherSys))
	MetricsValues.Store("PauseTotalNs", float64(m.MemStats.PauseTotalNs))
	MetricsValues.Store("StackInuse", float64(m.MemStats.StackInuse))
	MetricsValues.Store("StackSys", float64(m.MemStats.StackSys))
	MetricsValues.Store("Sys", float64(m.MemStats.Sys))
	MetricsValues.Store("TotalAlloc", float64(m.MemStats.TotalAlloc))
	MetricsValues.Store("RandomValue", random.Float64())

}

func (m MockMetrics) GetMetrics(MetricsValues *sync.Map) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	MetricsValues.Store("Alloc", m.MemStats["Alloc"])
	MetricsValues.Store("BuckHashSys", m.MemStats["BuckHashSys"])
	MetricsValues.Store("Frees", m.MemStats["Frees"])
	MetricsValues.Store("GCCPUFraction", m.MemStats["GCCPUFraction"])
	MetricsValues.Store("RandomValue", random.Float64())

}

var MetricNamesMap = map[string]struct{}{
	"Alloc":         {},
	"BuckHashSys":   {},
	"Frees":         {},
	"GCCPUFraction": {},
	"GCSys":         {},
	"HeapAlloc":     {},
	"HeapIdle":      {},
	"HeapInuse":     {},
	"HeapObjects":   {},
	"HeapReleased":  {},
	"HeapSys":       {},
	"LastGC":        {},
	"Lookups":       {},
	"MCacheInuse":   {},
	"MCacheSys":     {},
	"MSpanInuse":    {},
	"MSpanSys":      {},
	"Mallocs":       {},
	"NextGC":        {},
	"NumForcedGC":   {},
	"NumGC":         {},
	"OtherSys":      {},
	"PauseTotalNs":  {},
	"StackInuse":    {},
	"StackSys":      {},
	"Sys":           {},
	"TotalAlloc":    {},
	"RandomValue":   {},
}
