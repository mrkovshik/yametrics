package model

import "testing"

func TestMetrics_ValidateMetrics(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
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
				ID:    "test",
				MType: MetricTypeGauge,
				Delta: nil,
				Value: &testGauge,
			},
			true,
		},
		{"noName",
			fields{
				ID:    "",
				MType: MetricTypeGauge,
				Delta: nil,
				Value: &testGauge,
			},
			false,
		},
		{"noType",
			fields{
				ID:    "test",
				MType: "",
				Delta: nil,
				Value: &testGauge,
			},
			false,
		},
		{"wrongType",
			fields{
				ID:    "test",
				MType: "someWrongType",
				Delta: nil,
				Value: &testGauge,
			},
			false,
		},
		{"wrongValueType",
			fields{
				ID:    "test",
				MType: MetricTypeCounter,
				Delta: nil,
				Value: &testGauge,
			},
			false,
		},
		{"wrongValueType2",
			fields{
				ID:    "test",
				MType: MetricTypeGauge,
				Delta: &testCounter,
				Value: nil,
			},
			false,
		},
		{"nilValues",
			fields{
				ID:    "test",
				MType: MetricTypeCounter,
				Delta: nil,
				Value: nil,
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
			if got := m.ValidateMetrics(); got != tt.want {
				t.Errorf("ValidateMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
