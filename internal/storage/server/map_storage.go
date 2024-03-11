package server

import (
	"errors"
	"fmt"
)

type MapStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func NewMapStorage() *MapStorage {
	return &MapStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func (m *MapStorage) UpdateGauge(g Gauge) error {
	m.Gauges[g.name] = g.value
	return nil
}

func (m *MapStorage) UpdateCounter(c Counter) error {
	m.Counters[c.name] += c.value
	return nil
}

func (m *MapStorage) GetCounterValue(name string) (string, error) {
	value, ok := m.Counters[name]
	if !ok {
		return "", errors.New("not found")
	}
	return fmt.Sprint(value), nil
}

func (m *MapStorage) GetGaugeValue(name string) (string, error) {
	value, ok := m.Gauges[name]
	if !ok {
		return "", errors.New("not found")
	}
	return fmt.Sprint(value), nil
}

func (m *MapStorage) GetAllMetrics() string {

	resp := "<html><body><h1>Metric List</h1>" +
		"<h2>Gauges:</h2><ul>"

	for name, value := range m.Gauges {
		resp += fmt.Sprintf("<li><strong>%s:</strong> %f</li>", name, value)
	}
	resp += "</ul><h2>Counters:</h2><ul>"
	for name, value := range m.Counters {
		resp += fmt.Sprintf("<li><strong>%s:</strong> %v</li>", name, value)
	}
	resp += "</ul></body></html>"
	return resp
}
