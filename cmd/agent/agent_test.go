package main

import (
	"github.com/mrkovshik/yametrics/internal/metrics"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_getMetrics(t *testing.T) {
	var (
		src = metrics.NewMockMetrics()
		m   = sync.Map{}
	)
	tests := []struct {
		name string
	}{
		{"positive 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src.StoreMetrics(&m)
			valAlloc, ok := m.Load("Alloc")
			assert.True(t, ok)
			assert.Equal(t, valAlloc, 1.00)
			valBuckHashSys, ok := m.Load("BuckHashSys")
			assert.True(t, ok)
			assert.Equal(t, valBuckHashSys, 2.00)
		})
	}
}
