package storage

import (
	"fmt"
	"github.com/mrkovshik/yametrics/internal/metrics"
	"strconv"
)

type MemStorage struct {
	gauges   map[string]metrics.Gauge
	counters map[string]metrics.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]metrics.Gauge),
		counters: make(map[string]metrics.Counter),
	}
}

func (m *MemStorage) UpdateGauge(name, value string) error {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	m.gauges[name] = metrics.Gauge(intValue)
	fmt.Printf("gauge added\n name = %v\n, value = %v,\n MemStorage %v\n", name, value, m)
	return nil
}

func (m *MemStorage) UpdateCounter(name, value string) error {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	m.counters[name] += metrics.Counter(float64(intValue))
	fmt.Printf("counter added\n name = %v,\n value = %v,\n MemStorage %v\n", name, value, m)
	return nil
}
