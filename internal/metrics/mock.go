package metrics

import (
	"fmt"
	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
	"log"
	"math/rand"
	"time"
)

type MockMetrics struct {
	MemStats map[string]string
}

func NewMockMetrics() MockMetrics {
	return MockMetrics{
		map[string]string{
			"Alloc":         "1.00",
			"BuckHashSys":   "2.00",
			"Frees":         "3.00",
			"GCCPUFraction": "4.00",
		},
	}
}

func (m MockMetrics) PollMetrics(s storage.IAgentStorage) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	s.SaveMetric("Alloc", fmt.Sprint(m.MemStats["Alloc"]))
	s.SaveMetric("BuckHashSys", m.MemStats["BuckHashSys"])
	s.SaveMetric("Frees", m.MemStats["Frees"])
	s.SaveMetric("GCCPUFraction", m.MemStats["GCCPUFraction"])
	s.SaveMetric("RandomValue", fmt.Sprint(random.Float64()))
	if err := s.UpdateCounter(); err != nil {
		log.Fatal(err)
	}
}
