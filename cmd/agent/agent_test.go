package main

import (
	"github.com/mrkovshik/yametrics/internal/metrics"
	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getMetrics(t *testing.T) {
	var (
		src  = metrics.NewMockMetrics()
		strg = storage.NewAgentMapStorage()
	)
	tests := []struct {
		name string
	}{
		{"positive 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src.PollMetrics(strg)
			valPollCount := strg.LoadMetric("PollCount")
			assert.Equal(t, valPollCount, "1")
			valAlloc := strg.LoadMetric("Alloc")
			assert.Equal(t, valAlloc, "1.00")
			valBuckHashSys := strg.LoadMetric("BuckHashSys")
			assert.Equal(t, valBuckHashSys, "2.00")
			src.PollMetrics(strg)
			valPollCount = strg.LoadMetric("PollCount")
			assert.Equal(t, valPollCount, "2")
		})
	}
}
