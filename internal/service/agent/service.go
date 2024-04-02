package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mrkovshik/yametrics/internal/model"

	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	storage "github.com/mrkovshik/yametrics/internal/storage"
)

type Agent struct {
	Source  metrics.MetricSource
	Logger  *zap.SugaredLogger
	Config  config.AgentConfig
	Storage storage.IStorage
}

func NewAgent(source metrics.MetricSource, cfg config.AgentConfig, strg storage.IStorage, logger *zap.SugaredLogger) *Agent {
	return &Agent{
		Source:  source,
		Logger:  logger,
		Config:  cfg,
		Storage: strg,
	}
}

func (a *Agent) SendMetrics() {
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
	//a.Logger.Debug("Starting to send metrics")
	for {
		time.Sleep(time.Duration(a.Config.ReportInterval) * time.Second)
		for name := range metricNamesMap {
			go a.sendMetric(name)
		}
	}

}

func (a *Agent) PollMetrics() {

	for {
		//a.Logger.Debug("Starting to update metrics")
		a.Source.PollMetrics(a.Storage)
		time.Sleep(time.Duration(a.Config.PollInterval) * time.Second)
	}
}

func (a *Agent) sendMetric(name string) {
	var client = http.Client{Timeout: 30 * time.Second}
	currentMetric := model.Metrics{
		ID: name,
	}
	if name == "PollCount" {
		currentMetric.MType = model.MetricTypeCounter
	} else {
		currentMetric.MType = model.MetricTypeGauge
	}
	foundMetric, err := a.Storage.GetMetricByModel(currentMetric)
	if err != nil {
		a.Logger.Error("GetMetricByModel", err)
		return
	}
	metricUpdateURL := fmt.Sprintf("http://%v/update/", a.Config.Address)

	reqBuilder := NewRequestBuilder().SetURL(metricUpdateURL).AddJSONBody(foundMetric).Compress().SetMethod(http.MethodPost)
	if reqBuilder.Err != nil {
		a.Logger.Errorf("error building request: %v\nmetric name: %v", reqBuilder.Err, currentMetric.ID)
		return
	}
	response, err := client.Do(&reqBuilder.R)
	if err != nil {
		a.Logger.Errorf("error sending request: %v\nmetric name: %v", err, currentMetric.ID)
		return
	}
	if response.StatusCode != http.StatusOK {
		a.Logger.Errorf("status code is %v, while sending %v\n", response.StatusCode, currentMetric)
		return
	}
	if err := response.Body.Close(); err != nil {
		a.Logger.Error("response.Body.Close()", err)
		return
	}
}
