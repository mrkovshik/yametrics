package metrics

import (
	"context"
	"runtime"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/shirou/gopsutil/v3/mem"
)

type UtilMetrics struct {
	MemStats *mem.VirtualMemoryStat
}

func NewUtilMetrics() (UtilMetrics, error) {
	v, err := mem.VirtualMemory()
	return UtilMetrics{
		v,
	}, err
}

func (m UtilMetrics) PollMetrics(s service.Storage) error {
	ctx := context.Background()
	valueTotalMemory := float64(m.MemStats.Total)
	if err := s.UpdateMetricValue(ctx, model.Metrics{
		ID:    "TotalMemory",
		MType: model.MetricTypeGauge,
		Value: &valueTotalMemory,
	}); err != nil {
		return err
	}
	valueFreeMemory := float64(m.MemStats.Free)
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
