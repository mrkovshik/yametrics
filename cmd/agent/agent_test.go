package main

import (
	"testing"

	"github.com/mrkovshik/yametrics/internal/metrics"
	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
	"github.com/stretchr/testify/assert"
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
			valPollCount, err1 := strg.LoadCounter()
			assert.NoError(t, err1)
			assert.Equal(t, int64(1), valPollCount)
			valAlloc, err2 := strg.LoadMetric("Alloc")
			assert.NoError(t, err2)
			assert.Equal(t, 1.00, valAlloc)
			valBuckHashSys, err3 := strg.LoadMetric("BuckHashSys")
			assert.NoError(t, err3)
			assert.Equal(t, 2.00, valBuckHashSys)
			src.PollMetrics(strg)
			valPollCount, err1 = strg.LoadCounter()
			assert.NoError(t, err1)
			assert.Equal(t, int64(2), valPollCount)
		})
	}
}
