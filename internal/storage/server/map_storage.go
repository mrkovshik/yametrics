package storage

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/mrkovshik/yametrics/internal/metrics"
	"github.com/mrkovshik/yametrics/internal/templates"
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

func (s *MapStorage) UpdateMetricValue(metricType, metricName, metricValue string) error {
	switch metricType {

	case metrics.MetricTypeGauge:
		floatValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("wrong value format")
		}
		s.Gauges[metricName] = floatValue
	case metrics.MetricTypeCounter:
		intValue, err := strconv.ParseInt(metricValue, 0, 64)
		if err != nil {
			return errors.New("wrong value format")
		}
		s.Counters[metricName] += intValue
	default:
		return errors.New("invalid metric type")
	}
	return nil
}

func (s *MapStorage) GetMetricValue(metricType, metricName string) (string, error) {
	var stringValue string
	switch metricType {
	case metrics.MetricTypeGauge:
		value, ok := s.Gauges[metricName]
		if !ok {
			return "", errors.New("not found")
		}
		stringValue = fmt.Sprint(value)
	case metrics.MetricTypeCounter:
		value, ok := s.Counters[metricName]
		if !ok {
			return "", errors.New("not found")
		}
		stringValue = fmt.Sprint(value)
	}
	return stringValue, nil
}

func (s *MapStorage) GetAllMetrics() (string, error) {
	var tpl bytes.Buffer
	t, err := templates.ParseTemplates()
	if err != nil {
		return "", err
	}
	if err := t.ExecuteTemplate(&tpl, "list_metrics", s); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
