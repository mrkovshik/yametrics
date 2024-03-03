package storage

import "strconv"

type (
	gauge struct {
		name  string
		value float64
	}
)

func (g gauge) Update(s IStorage) error {
	return s.UpdateGauge(g)
}

func NewGauge(name, value string) (gauge, error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return gauge{}, err
	}
	return gauge{name: name, value: floatValue}, err
}
