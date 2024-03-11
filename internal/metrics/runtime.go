package metrics

import (
	"fmt"
	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
	"log"
	"math/rand"
	"runtime"
	"time"
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
	s.SaveMetric("Alloc", fmt.Sprint(m.MemStats.Alloc))
	s.SaveMetric("Alloc", fmt.Sprint(m.MemStats.Alloc))
	s.SaveMetric("BuckHashSys", fmt.Sprint(m.MemStats.BuckHashSys))
	s.SaveMetric("Frees", fmt.Sprint(m.MemStats.Frees))
	s.SaveMetric("GCCPUFraction", fmt.Sprint(m.MemStats.GCCPUFraction))
	s.SaveMetric("GCSys", fmt.Sprint(m.MemStats.GCSys))
	s.SaveMetric("HeapAlloc", fmt.Sprint(m.MemStats.HeapAlloc))
	s.SaveMetric("HeapIdle", fmt.Sprint(m.MemStats.HeapIdle))
	s.SaveMetric("HeapInuse", fmt.Sprint(m.MemStats.HeapInuse))
	s.SaveMetric("HeapObjects", fmt.Sprint(m.MemStats.HeapObjects))
	s.SaveMetric("HeapReleased", fmt.Sprint(m.MemStats.HeapReleased))
	s.SaveMetric("HeapSys", fmt.Sprint(m.MemStats.HeapSys))
	s.SaveMetric("LastGC", fmt.Sprint(m.MemStats.LastGC))
	s.SaveMetric("Lookups", fmt.Sprint(m.MemStats.Lookups))
	s.SaveMetric("MCacheInuse", fmt.Sprint(m.MemStats.MCacheInuse))
	s.SaveMetric("MCacheSys", fmt.Sprint(m.MemStats.MCacheSys))
	s.SaveMetric("MSpanInuse", fmt.Sprint(m.MemStats.MSpanInuse))
	s.SaveMetric("MSpanSys", fmt.Sprint(m.MemStats.MSpanSys))
	s.SaveMetric("Mallocs", fmt.Sprint(m.MemStats.Mallocs))
	s.SaveMetric("NextGC", fmt.Sprint(m.MemStats.NextGC))
	s.SaveMetric("NumForcedGC", fmt.Sprint(m.MemStats.NumForcedGC))
	s.SaveMetric("NumGC", fmt.Sprint(m.MemStats.NumGC))
	s.SaveMetric("OtherSys", fmt.Sprint(m.MemStats.OtherSys))
	s.SaveMetric("PauseTotalNs", fmt.Sprint(m.MemStats.PauseTotalNs))
	s.SaveMetric("StackInuse", fmt.Sprint(m.MemStats.StackInuse))
	s.SaveMetric("StackSys", fmt.Sprint(m.MemStats.StackSys))
	s.SaveMetric("Sys", fmt.Sprint(m.MemStats.Sys))
	s.SaveMetric("TotalAlloc", fmt.Sprint(m.MemStats.TotalAlloc))
	s.SaveMetric("RandomValue", fmt.Sprint(random.Float64()))
	if err := s.UpdateCounter(); err != nil {
		log.Fatal(err)
	}
}
