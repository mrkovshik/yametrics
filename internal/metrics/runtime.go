package metrics

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
)

type RuntimeMetrics struct {
	MemStats runtime.MemStats
}

func NewRuntimeMetrics() RuntimeMetrics {
	m := RuntimeMetrics{
		MemStats: runtime.MemStats{},
	}
	runtime.ReadMemStats(&m.MemStats)
	return m
}

func (m RuntimeMetrics) PollMetrics(s storage.IAgentStorage) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.SaveMetric("Alloc", float64(m.MemStats.Alloc))
	s.SaveMetric("BuckHashSys", float64(m.MemStats.BuckHashSys))
	s.SaveMetric("Frees", float64(m.MemStats.Frees))
	s.SaveMetric("GCCPUFraction", m.MemStats.GCCPUFraction)
	s.SaveMetric("GCSys", float64(m.MemStats.GCSys))
	s.SaveMetric("HeapAlloc", float64(m.MemStats.HeapAlloc))
	s.SaveMetric("HeapIdle", float64(m.MemStats.HeapIdle))
	s.SaveMetric("HeapInuse", float64(m.MemStats.HeapInuse))
	s.SaveMetric("HeapObjects", float64(m.MemStats.HeapObjects))
	s.SaveMetric("HeapReleased", float64(m.MemStats.HeapReleased))
	s.SaveMetric("HeapSys", float64(m.MemStats.HeapSys))
	s.SaveMetric("LastGC", float64(m.MemStats.LastGC))
	s.SaveMetric("Lookups", float64(m.MemStats.Lookups))
	s.SaveMetric("MCacheInuse", float64(m.MemStats.MCacheInuse))
	s.SaveMetric("MCacheSys", float64(m.MemStats.MCacheSys))
	s.SaveMetric("MSpanInuse", float64(m.MemStats.MSpanInuse))
	s.SaveMetric("MSpanSys", float64(m.MemStats.MSpanSys))
	s.SaveMetric("Mallocs", float64(m.MemStats.Mallocs))
	s.SaveMetric("NextGC", float64(m.MemStats.NextGC))
	s.SaveMetric("NumForcedGC", float64(m.MemStats.NumForcedGC))
	s.SaveMetric("NumGC", float64(m.MemStats.NumGC))
	s.SaveMetric("OtherSys", float64(m.MemStats.OtherSys))
	s.SaveMetric("PauseTotalNs", float64(m.MemStats.PauseTotalNs))
	s.SaveMetric("StackInuse", float64(m.MemStats.StackInuse))
	s.SaveMetric("StackSys", float64(m.MemStats.StackSys))
	s.SaveMetric("Sys", float64(m.MemStats.Sys))
	s.SaveMetric("TotalAlloc", float64(m.MemStats.TotalAlloc))
	s.SaveMetric("RandomValue", random.Float64())
	if err := s.UpdateCounter(); err != nil {
		log.Fatal(err)
	}
}
