// Package metrics provides functionality to collect runtime and virtual memory metrics.
package metrics

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/mem"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/service"
)

// RuntimeMetrics holds memory statistics from the runtime package and virtual memory statistics.
type RuntimeMetrics struct {
	MemStats     runtime.MemStats       // Statistics on the memory allocator.
	VirtMemStats *mem.VirtualMemoryStat // Virtual memory statistics.
}

// NewRuntimeMetrics creates a new RuntimeMetrics instance with initialized memory statistics.
func NewRuntimeMetrics() RuntimeMetrics {
	m := RuntimeMetrics{
		MemStats:     runtime.MemStats{},
		VirtMemStats: &mem.VirtualMemoryStat{},
	}
	return m
}

// PollMemStats collects memory statistics and updates them in the provided storage service.
func (m RuntimeMetrics) PollMemStats(s service.Storage) error {
	ctx := context.Background()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	runtime.ReadMemStats(&m.MemStats)
	valueAlloc := float64(m.MemStats.Alloc)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "Alloc",
		MType: model.MetricTypeGauge,
		Value: &valueAlloc,
	}); err != nil {
		return err
	}
	valueBuckHashSys := float64(m.MemStats.BuckHashSys)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "BuckHashSys",
		MType: model.MetricTypeGauge,
		Value: &valueBuckHashSys,
	}); err != nil {
		return err
	}
	valueFrees := float64(m.MemStats.Frees)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "Frees",
		MType: model.MetricTypeGauge,
		Value: &valueFrees,
	}); err != nil {
		return err
	}
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "GCCPUFraction",
		MType: model.MetricTypeGauge,
		Value: &m.MemStats.GCCPUFraction,
	}); err != nil {
		return err
	}
	valueGCSys := float64(m.MemStats.GCSys)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "GCSys",
		MType: model.MetricTypeGauge,
		Value: &valueGCSys,
	}); err != nil {
		return err
	}
	valueHeapAlloc := float64(m.MemStats.HeapAlloc)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "HeapAlloc",
		MType: model.MetricTypeGauge,
		Value: &valueHeapAlloc,
	}); err != nil {
		return err
	}
	valueHeapIdle := float64(m.MemStats.HeapIdle)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "HeapIdle",
		MType: model.MetricTypeGauge,
		Value: &valueHeapIdle,
	}); err != nil {
		return err
	}
	valueHeapInuse := float64(m.MemStats.HeapInuse)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "HeapInuse",
		MType: model.MetricTypeGauge,
		Value: &valueHeapInuse,
	}); err != nil {
		return err
	}
	valueHeapObjects := float64(m.MemStats.HeapObjects)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "HeapObjects",
		MType: model.MetricTypeGauge,
		Value: &valueHeapObjects,
	}); err != nil {
		return err
	}
	valueHeapReleased := float64(m.MemStats.HeapReleased)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "HeapReleased",
		MType: model.MetricTypeGauge,
		Value: &valueHeapReleased,
	}); err != nil {
		return err
	}
	valueHeapSys := float64(m.MemStats.HeapSys)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "HeapSys",
		MType: model.MetricTypeGauge,
		Value: &valueHeapSys,
	}); err != nil {
		return err
	}
	valueLastGC := float64(m.MemStats.LastGC)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "LastGC",
		MType: model.MetricTypeGauge,
		Value: &valueLastGC,
	}); err != nil {
		return err
	}
	valueLookups := float64(m.MemStats.Lookups)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "Lookups",
		MType: model.MetricTypeGauge,
		Value: &valueLookups,
	}); err != nil {
		return err
	}
	valueMCacheInuse := float64(m.MemStats.MCacheInuse)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "MCacheInuse",
		MType: model.MetricTypeGauge,
		Value: &valueMCacheInuse,
	}); err != nil {
		return err
	}
	valueMCacheSys := float64(m.MemStats.MCacheSys)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "MCacheSys",
		MType: model.MetricTypeGauge,
		Value: &valueMCacheSys,
	}); err != nil {
		return err
	}
	valueMSpanInuse := float64(m.MemStats.MSpanInuse)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "MSpanInuse",
		MType: model.MetricTypeGauge,
		Value: &valueMSpanInuse,
	}); err != nil {
		return err
	}
	valueMSpanSys := float64(m.MemStats.MSpanSys)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "MSpanSys",
		MType: model.MetricTypeGauge,
		Value: &valueMSpanSys,
	}); err != nil {
		return err
	}
	valueMallocs := float64(m.MemStats.Mallocs)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "Mallocs",
		MType: model.MetricTypeGauge,
		Value: &valueMallocs,
	}); err != nil {
		return err
	}
	valueNextGC := float64(m.MemStats.NextGC)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "NextGC",
		MType: model.MetricTypeGauge,
		Value: &valueNextGC,
	}); err != nil {
		return err
	}
	valueNumForcedGC := float64(m.MemStats.NumForcedGC)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "NumForcedGC",
		MType: model.MetricTypeGauge,
		Value: &valueNumForcedGC,
	}); err != nil {
		return err
	}
	valueNumGC := float64(m.MemStats.NumGC)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "NumGC",
		MType: model.MetricTypeGauge,
		Value: &valueNumGC,
	}); err != nil {
		return err
	}
	valueOtherSys := float64(m.MemStats.OtherSys)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "OtherSys",
		MType: model.MetricTypeGauge,
		Value: &valueOtherSys,
	}); err != nil {
		return err
	}
	valuePauseTotalNs := float64(m.MemStats.PauseTotalNs)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "PauseTotalNs",
		MType: model.MetricTypeGauge,
		Value: &valuePauseTotalNs,
	}); err != nil {
		return err
	}
	valueStackInuse := float64(m.MemStats.StackInuse)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "StackInuse",
		MType: model.MetricTypeGauge,
		Value: &valueStackInuse,
	}); err != nil {
		return err
	}
	valueStackSys := float64(m.MemStats.StackSys)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "StackSys",
		MType: model.MetricTypeGauge,
		Value: &valueStackSys,
	}); err != nil {
		return err
	}
	valueSys := float64(m.MemStats.Sys)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "Sys",
		MType: model.MetricTypeGauge,
		Value: &valueSys,
	}); err != nil {
		return err
	}
	valueTotalAlloc := float64(m.MemStats.TotalAlloc)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "TotalAlloc",
		MType: model.MetricTypeGauge,
		Value: &valueTotalAlloc,
	}); err != nil {
		return err
	}
	valueRandom := random.Float64()
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "RandomValue",
		MType: model.MetricTypeGauge,
		Value: &valueRandom,
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

// PollVirtMemStats collects virtual memory statistics and updates them in the provided storage service.
func (m RuntimeMetrics) PollVirtMemStats(s service.Storage) error {
	ctx := context.Background()
	var err error
	m.VirtMemStats, err = mem.VirtualMemory()
	if err != nil {
		return err
	}
	valueTotalMemory := float64(m.VirtMemStats.Total)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "TotalMemory",
		MType: model.MetricTypeGauge,
		Value: &valueTotalMemory,
	}); err != nil {
		return err
	}
	valueFreeMemory := float64(m.VirtMemStats.Free)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "FreeMemory",
		MType: model.MetricTypeGauge,
		Value: &valueFreeMemory,
	}); err != nil {
		return err
	}
	valueCPUutilization1 := float64(runtime.NumCPU())
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "CPUutilization1",
		MType: model.MetricTypeGauge,
		Value: &valueCPUutilization1,
	}); err != nil {
		return err
	}
	return nil
}
