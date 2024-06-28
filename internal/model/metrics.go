package model

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Constants for metric types.
const (
	// MetricTypeGauge represents a gauge metric type.
	MetricTypeGauge = "gauge"

	// MetricTypeCounter represents a counter metric type.
	MetricTypeCounter = "counter"
)

// Metrics represents a metric entity with ID, type (gauge or counter), and either Delta (for counter) or Value (for gauge).
type Metrics struct {
	ID    string   `json:"id"`              // Metric name
	MType string   `json:"type"`            // Parameter that takes values gauge or counter
	Delta *int64   `json:"delta,omitempty"` // Metric value in case of counter transmission
	Value *float64 `json:"value,omitempty"` // Metric value in case of gauge transmission
}

// MapMetricsFromReqJSON maps metric data from JSON format in the HTTP request body to Metrics struct.
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

// MapMetricsFromReqURL maps metric data from URL parameters to Metrics struct.
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
