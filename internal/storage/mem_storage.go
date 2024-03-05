package storage

import (
	"fmt"
)

type (
	MapStorage struct {
		Gauges   map[string]float64
		Counters map[string]int64
	}
	IStorage interface {
		UpdateCounter(counter) error
		UpdateGauge(gauge) error
	}
)

func NewMapStorage() *MapStorage {
	return &MapStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func (m *MapStorage) UpdateGauge(g gauge) error {
	m.Gauges[g.name] = g.value
	fmt.Printf("gauge added\n name = %v\n, value = %v,\n MemStorage %v\n", g.name, g.value, m)
	return nil
}

func (m *MapStorage) UpdateCounter(c counter) error {
	m.Counters[c.name] += c.value
	fmt.Printf("gauge added\n name = %v\n, value = %v,\n MemStorage %v\n", c.name, c.value, m)
	return nil
}
