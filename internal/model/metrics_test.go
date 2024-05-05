package model

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_ValidateMetrics(t *testing.T) {
	type fields struct {
		ID     string
		MType  string
		Delta  *int64
		Value  *float64
		Method string
	}
	var (
		testGauge         = 3.6
		testCounter int64 = 3
	)
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"pos 1",
			fields{
				ID:     "test",
				MType:  MetricTypeGauge,
				Delta:  nil,
				Value:  &testGauge,
				Method: http.MethodPost,
			},
			true,
		},
		{"noName",
			fields{
				ID:     "",
				MType:  MetricTypeGauge,
				Delta:  nil,
				Value:  &testGauge,
				Method: http.MethodPost,
			},
			false,
		},
		{"noType",
			fields{
				ID:     "test",
				MType:  "",
				Delta:  nil,
				Value:  &testGauge,
				Method: http.MethodPost,
			},
			false,
		},
		{"wrongType",
			fields{
				ID:     "test",
				MType:  "someWrongType",
				Delta:  nil,
				Value:  &testGauge,
				Method: http.MethodPost,
			},
			false,
		},
		{"wrongValueType",
			fields{
				ID:     "test",
				MType:  MetricTypeCounter,
				Delta:  nil,
				Value:  &testGauge,
				Method: http.MethodPost,
			},
			false,
		},
		{"wrongValueType2",
			fields{
				ID:     "test",
				MType:  MetricTypeGauge,
				Delta:  &testCounter,
				Value:  nil,
				Method: http.MethodPost,
			},
			false,
		},
		{"nilValues",
			fields{
				ID:     "test",
				MType:  MetricTypeCounter,
				Delta:  nil,
				Value:  nil,
				Method: http.MethodPost,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Metrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			buf := bytes.Buffer{}
			err1 := json.NewEncoder(&buf).Encode(m)
			assert.NoError(t, err1)
			req, err2 := http.NewRequest(tt.fields.Method, "test", &buf)
			assert.NoError(t, err2)

			err3 := m.MapMetricsFromReqJSON(req)
			if (err3 == nil) != tt.want {
				t.Errorf("MapMetricsFromReqJSON() = %v, want %v", err3 == nil, tt.want)
			}
		})
	}
}
