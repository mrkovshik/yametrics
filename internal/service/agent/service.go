package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mrkovshik/yametrics/internal/model"
	"net/http"
	"time"

	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
)

type Agent struct {
	Source  metrics.MetricSource
	Logger  *zap.SugaredLogger
	Config  config.AgentConfig
	Storage storage.IAgentStorage
}

func NewAgent(source metrics.MetricSource, cfg config.AgentConfig, strg storage.IAgentStorage, logger *zap.SugaredLogger) *Agent {
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
	currentMetric := model.Metrics{
		ID: name,
	}
	if name == "PollCount" {
		delta, err := a.Storage.LoadCounter()
		if err != nil {
			a.Logger.Error("a.Storage.LoadCounter", err)
			return
		}

		currentMetric.MType = model.MetricTypeCounter
		currentMetric.Delta = &delta
	} else {
		value, err := a.Storage.LoadMetric(name)
		if err != nil {
			a.Logger.Error("a.Storage.LoadMetric", err)
			return
		}
		currentMetric.MType = model.MetricTypeGauge
		currentMetric.Value = &value
	}

	metricUpdateURL := fmt.Sprintf("http://%v/update/", a.Config.Address)
	buf := bytes.Buffer{}
	if err3 := json.NewEncoder(&buf).Encode(currentMetric); err3 != nil {
		a.Logger.Error("Encode", zap.Error(err3))
		return
	}

	response, err := http.Post(metricUpdateURL, "application/json", &buf)
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
