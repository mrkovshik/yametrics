package storage

import (
	"fmt"
)

func NewMapStorage() *MapStorage {
	return &MapStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

func (m *MapStorage) UpdateGauge(g Gauge) error {
	m.Gauges[g.name] = g.value
	fmt.Printf("Gauge added\n name = %v\n, value = %v,\n MemStorage %v\n", g.name, g.value, m)
	return nil
}

func (m *MapStorage) UpdateCounter(c Counter) error {
	m.Counters[c.name] += c.value
	fmt.Printf("Gauge added\n name = %v\n, value = %v,\n MemStorage %v\n", c.name, c.value, m)
	return nil
}

func (m *MapStorage) GetCounterValue(name string) string {
	return fmt.Sprint(m.Counters[name])
}

func (m *MapStorage) GetGaugeValue(name string) string {
	return fmt.Sprint(m.Gauges[name])
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
