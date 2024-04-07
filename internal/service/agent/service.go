package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mrkovshik/yametrics/internal/model"
	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	"github.com/mrkovshik/yametrics/internal/storage"
)

type Agent struct {
	source  metrics.MetricSource
	logger  *zap.SugaredLogger
	config  config.AgentConfig
	storage storage.Storage
}

func NewAgent(source metrics.MetricSource, cfg config.AgentConfig, strg storage.Storage, logger *zap.SugaredLogger) *Agent {
	return &Agent{
		source:  source,
		logger:  logger,
		config:  cfg,
		storage: strg,
	}
}

func (a *Agent) SendMetrics(ctx context.Context) {
	var metricNamesMap = map[string]struct{}{
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
	//a.logger.Debug("Starting to send metrics")
	for {
		time.Sleep(time.Duration(a.config.ReportInterval) * time.Second)
		for name := range metricNamesMap {
			go a.sendMetric(ctx, name)
		}
	}

}

func (a *Agent) PollMetrics() {

	for {
		//a.logger.Debug("Starting to update metrics")
		a.source.PollMetrics(a.storage)
		time.Sleep(time.Duration(a.config.PollInterval) * time.Second)
	}
}

func (a *Agent) sendMetric(ctx context.Context, name string) {
	var client = http.Client{Timeout: 30 * time.Second}
	currentMetric := model.Metrics{
		ID: name,
	}
	if name == "PollCount" {
		currentMetric.MType = model.MetricTypeCounter
	} else {
		currentMetric.MType = model.MetricTypeGauge
	}
	foundMetric, err := a.storage.GetMetricByModel(ctx, currentMetric)
	if err != nil {
		a.logger.Error("GetMetricByModel", err)
		return
	}
	metricUpdateURL := fmt.Sprintf("http://%v/update/", a.config.Address)

	reqBuilder := NewRequestBuilder().SetURL(metricUpdateURL).AddJSONBody(foundMetric).Compress().SetMethod(http.MethodPost)
	if reqBuilder.Err != nil {
		a.logger.Errorf("error building request: %v\nmetric name: %v", reqBuilder.Err, currentMetric.ID)
		return
	}
	response, err := client.Do(&reqBuilder.R)
	if err != nil {
		a.logger.Errorf("error sending request: %v\nmetric name: %v", err, currentMetric.ID)
		return
	}
	if response.StatusCode != http.StatusOK {
		a.logger.Errorf("status code is %v, while sending %v\n", response.StatusCode, currentMetric)
		return
	}
	if err := response.Body.Close(); err != nil {
		a.logger.Error("response.Body.Close()", err)
		return
	}
}
