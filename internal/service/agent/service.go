package service

import (
	"errors"
	"fmt"
	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	storage "github.com/mrkovshik/yametrics/internal/storage/agent"
	"log"
	"net/http"
	"time"
)

type Agent struct {
	Source  metrics.MetricSource
	Logger  *log.Logger
	Config  config.AgentConfig
	Storage storage.IAgentStorage
}

func NewAgent(source metrics.MetricSource, logger *log.Logger, cfg config.AgentConfig, strg storage.IAgentStorage) *Agent {
	return &Agent{
		Source:  source,
		Logger:  logger,
		Config:  cfg,
		Storage: strg,
	}
}

func (a *Agent) SendMetric() error {

	log.Println("Starting to send metrics")
	for name := range metrics.MetricNamesMap {
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
			log.Printf("status code is %v, while sending %v:%v:%v\n", response.StatusCode, metricType, name, value)
			return errors.New("status code is not OK")
		}
		response.Body.Close()
	}

	return nil
}

func (a *Agent) PollMetrics() {

	for {
		log.Println("Starting to update metrics")
		a.Source.PollMetrics(a.Storage)

		time.Sleep(time.Duration(a.Config.PollInterval) * time.Second)
	}
}
