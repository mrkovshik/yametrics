package service

import (
	"bytes"
	"encoding/json"
	"errors"
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
		currentMetric := model.Metrics{
			ID: name,
		}
		if name == "PollCount" {

			delta, err := a.Storage.LoadCounter()
			if err != nil {
				a.Logger.Error("a.Storage.LoadCounter", err)
				return err
			}

			currentMetric.MType = model.MetricTypeCounter
			currentMetric.Delta = &delta
		} else {

			delta, err := a.Storage.LoadMetric(name)
			if err != nil {
				a.Logger.Error("a.Storage.LoadMetric", err)
				return err
			}
			currentMetric.MType = model.MetricTypeGauge
			currentMetric.Value = &delta
		}

		metricUpdateURL := fmt.Sprintf("http://%v/update/", a.Config.Address)
		buf := bytes.Buffer{}
		if err3 := json.NewEncoder(&buf).Encode(currentMetric); err3 != nil {
			a.Logger.Error("Encode", zap.Error(err3))
			return err3
		}

		response, err := http.Post(metricUpdateURL, "application/json", &buf)
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			a.Logger.Errorf("status code is %v, while sending %v\n", response.StatusCode, currentMetric)
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
