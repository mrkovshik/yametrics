package model

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

func (m Metrics) ValidateMetrics() bool {
	return m.ID != "" && (m.MType == MetricTypeGauge && (m.Delta == nil && m.Value != nil)) || (m.MType == MetricTypeCounter && (m.Delta != nil && m.Value == nil))
}
