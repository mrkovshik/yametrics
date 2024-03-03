package storage

import (
	"fmt"
)

type (
	MapStorage struct {
		gauges   map[string]float64
		counters map[string]int64
	}
	IStorage interface {
		UpdateCounter(counter) error
		UpdateGauge(gauge) error
	}
)

func NewMapStorage() *MapStorage {
	return &MapStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MapStorage) UpdateGauge(g gauge) error {
	m.gauges[g.name] = g.value
	fmt.Printf("gauge added\n name = %v\n, value = %v,\n MemStorage %v\n", g.name, g.value, m)
	return nil
}

func (m *MapStorage) UpdateCounter(c counter) error {
	m.counters[c.name] += c.value
	fmt.Printf("gauge added\n name = %v\n, value = %v,\n MemStorage %v\n", c.name, c.value, m)
	return nil
}
