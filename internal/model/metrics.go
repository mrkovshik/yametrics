package model

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metrics) MapMetricsFromReqJSON(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		return err
	}
	switch m.MType {
	case MetricTypeGauge:
		if req.Method == http.MethodPost {
			if m.Value == nil {
				return errors.New("errInvalidMetricType")
			}
		}
	case MetricTypeCounter:
		if req.Method == http.MethodPost {
			if m.Delta == nil {
				return errors.New("errInvalidMetricType")
			}
		}
	default:
		return errors.New("errInvalidMetricType")
	}
	if m.ID == "" {
		return errors.New("errInvalidMetricType")
	}
	return nil

}

func (m *Metrics) MapMetricsFromReqURL(req *http.Request) error {
	metricName := chi.URLParam(req, "name")
	metricValue := chi.URLParam(req, "value")
	metricType := chi.URLParam(req, "type")

	switch metricType {
	case MetricTypeGauge:
		if req.Method == http.MethodPost {
			floatVal, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				return err
			}
			m.Value = &floatVal
		}
	case MetricTypeCounter:
		if req.Method == http.MethodPost {
			intVal, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				return err
			}
			m.Delta = &intVal
		}
	default:
		return errors.New("errInvalidMetricType")
	}
	m.ID = metricName
	m.MType = metricType
	return nil
}
