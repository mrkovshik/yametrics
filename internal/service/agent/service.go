package service

import (
	"errors"
	"fmt"
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

func (a *Agent) SendMetric() error {
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
	a.Logger.Debug("Starting to send metrics")
	for name := range metricNamesMap {
		metricType := metrics.MetricTypeGauge
		if name == "PollCount" {
			metricType = metrics.MetricTypeCounter
		}
		value := a.Storage.LoadMetric(name)
		metricUpdateURL := fmt.Sprintf("http://%v/update/%v/%v/%v", a.Config.Address, metricType, name, value)
		response, err := http.Post(metricUpdateURL, "text/plain", nil)
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			a.Logger.Errorf("status code is %v, while sending %v:%v:%v\n", response.StatusCode, metricType, name, value)
			return errors.New("status code is not OK")
		}
		if err := response.Body.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (a *Agent) PollMetrics() {

	for {
		a.Logger.Debug("Starting to update metrics")
		a.Source.PollMetrics(a.Storage)
		time.Sleep(time.Duration(a.Config.PollInterval) * time.Second)
	}
}
