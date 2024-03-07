package storage

type (
	Gauge struct {
		name  string
		value float64
	}
)

func (g Gauge) Update(s IStorage) error {
	return s.UpdateGauge(g)
}

func NewGauge(name string, value float64) Gauge {
	return Gauge{name: name, value: value}
}
