package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mrkovshik/yametrics/internal/storage"
)

func TestRuntimeMetrics_PollMemStats(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		m := NewRuntimeMetrics()
		s := storage.NewMapStorage()
		err := m.PollMemStats(s)
		assert.NoError(t, err)
	})
}

func TestRuntimeMetrics_PollVirtMemStats(t *testing.T) {

	t.Run("1", func(t *testing.T) {
		m := NewRuntimeMetrics()
		s := storage.NewMapStorage()
		err := m.PollVirtMemStats(s)
		assert.NoError(t, err)
	})
}
