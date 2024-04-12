package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/service"
	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
)

type Agent struct {
	source  metrics.MetricSource
	logger  *zap.SugaredLogger
	config  config.AgentConfig
	storage service.Storage
}

func NewAgent(source metrics.MetricSource, cfg config.AgentConfig, strg service.Storage, logger *zap.SugaredLogger) *Agent {
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
		a.sendMetrics(ctx, metricNamesMap)
	}

}

func (a *Agent) PollMetrics() {

	for {
		//a.logger.Debug("Starting to update metrics")
		if err := a.source.PollMetrics(a.storage); err != nil {
			a.logger.Error("PollMetrics", err)
			return
		}
		time.Sleep(time.Duration(a.config.PollInterval) * time.Second)
	}
}

func (a *Agent) sendMetrics(ctx context.Context, names map[string]struct{}) {
	var batch []model.Metrics

	for name := range names {
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
		batch = append(batch, foundMetric)
	}

	metricUpdateURL := fmt.Sprintf("http://%v/updates/", a.config.Address)

	reqBuilder := NewRequestBuilder().SetURL(metricUpdateURL).AddJSONBody(batch).Compress().SetMethod(http.MethodPost)
	if reqBuilder.Err != nil {
		a.logger.Errorf("error building request: %v\n", reqBuilder.Err)
		return
	}
	response, err := a.retryableSend(&reqBuilder.R)
	if err != nil {
		a.logger.Errorf("error sending request: %v\n", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		a.logger.Errorf("status code is %v\n", response.StatusCode)
		return
	}
	if err := response.Body.Close(); err != nil {
		a.logger.Error("response.Body.Close()", err)
		return
	}
}

func (a *Agent) retryableSend(req *http.Request) (*http.Response, error) {
	var (
		bodyBytes      []byte
		retryIntervals = []int{1, 3, 5} //TODO: move to config
		client         = http.Client{Timeout: 5 * time.Second}
		err            error
	)
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	for i := 0; i <= len(retryIntervals); i++ {
		response, err := client.Do(req)
		if err == nil {
			return response, nil
		}
		if i == len(retryIntervals) {
			return nil, err
		}
		a.logger.Errorf("failed connect to server: %v\n retry in %v seconds\n", err, retryIntervals[i])
		time.Sleep(time.Duration(retryIntervals[i]) * time.Second)
		if req.Body != nil {
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
	}
	return nil, nil
}
