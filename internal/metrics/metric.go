package metrics

import "github.com/mrkovshik/yametrics/internal/storage/server"

type Imetric interface { //TODO: Подумать, нужен ли тут интерфейс вообще
	Update(server.IStorage) error
}

const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

var MetricNamesMap = map[string]struct{}{
	"Alloc":         {},
	"BuckHashSys":   {},
	"Frees":         {},
	"GCCPUFraction": {},
	"GCSys":         {},
	"HeapAlloc":     {},
	"HeapIdle":      {},
	"HeapInuse":     {},
	"HeapObjects":   {},
	"HeapReleased":  {},
	"HeapSys":       {},
	"LastGC":        {},
	"Lookups":       {},
	"MCacheInuse":   {},
	"MCacheSys":     {},
	"MSpanInuse":    {},
	"MSpanSys":      {},
	"Mallocs":       {},
	"NextGC":        {},
	"NumForcedGC":   {},
	"NumGC":         {},
	"OtherSys":      {},
	"PauseTotalNs":  {},
	"StackInuse":    {},
	"StackSys":      {},
	"Sys":           {},
	"TotalAlloc":    {},
	"RandomValue":   {},
	"PollCount":     {},
}
