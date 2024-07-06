// Package service provides utilities for building and sending HTTP requests with metrics to server.
package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/model"
	"go.uber.org/zap"

	"github.com/mrkovshik/yametrics/internal/metrics"
)

type storage interface {
	UpdateMetricValue(ctx context.Context, newMetrics model.Metrics) error

	UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error

	GetMetricByModel(ctx context.Context, newMetrics model.Metrics) (model.Metrics, error)

	GetAllMetrics(ctx context.Context) (map[string]model.Metrics, error)
}

// Agent represents a metric collection agent that polls and sends metrics.
type Agent struct {
	source  metrics.MetricSource // Source of the metrics
	logger  *zap.SugaredLogger   // Logger for logging messages
	cfg     *config.AgentConfig  // Configuration for the agent
	storage storage              // Storage for metrics
}

// NewAgent initializes a new Agent.
func NewAgent(source metrics.MetricSource, cfg *config.AgentConfig, strg storage, logger *zap.SugaredLogger) *Agent {
	return &Agent{
		source:  source,
		logger:  logger,
		cfg:     cfg,
		storage: strg,
	}
}

// SendMetrics sends metrics at intervals specified by the channel.
func (a *Agent) SendMetrics(ctx context.Context, ch <-chan time.Time) {
	var metricNamesMap = map[string]struct{}{
		"Alloc":           {},
		"BuckHashSys":     {},
		"Frees":           {},
		"GCCPUFraction":   {},
		"GCSys":           {},
		"HeapAlloc":       {},
		"HeapIdle":        {},
		"HeapInuse":       {},
		"HeapObjects":     {},
		"HeapReleased":    {},
		"HeapSys":         {},
		"LastGC":          {},
		"Lookups":         {},
		"MCacheInuse":     {},
		"MCacheSys":       {},
		"MSpanInuse":      {},
		"MSpanSys":        {},
		"Mallocs":         {},
		"NextGC":          {},
		"NumForcedGC":     {},
		"NumGC":           {},
		"OtherSys":        {},
		"PauseTotalNs":    {},
		"StackInuse":      {},
		"StackSys":        {},
		"Sys":             {},
		"TotalAlloc":      {},
		"RandomValue":     {},
		"PollCount":       {},
		"TotalMemory":     {},
		"FreeMemory":      {},
		"CPUutilization1": {},
	}
	for range ch {
		a.sendMetricsByPool(ctx, metricNamesMap)
	}
}

// PollMetrics polls metrics at intervals specified by the channel.
func (a *Agent) PollMetrics(ch <-chan time.Time) {
	for range ch {
		a.logger.Debug("Starting to update metrics")
		if err := a.source.PollMemStats(a.storage); err != nil {
			a.logger.Error("PollMemStats", err)
			return
		}
	}
}

// PollUitlMetrics polls utilization metrics at intervals specified by the channel.
func (a *Agent) PollUitlMetrics(ch <-chan time.Time) {
	for range ch {
		if err := a.source.PollVirtMemStats(a.storage); err != nil {
			a.logger.Error("PollVirtMemStats", err)
			return
		}
	}
}

// sendMetricsByPool sends metrics using a pool of workers.
func (a *Agent) sendMetricsByPool(ctx context.Context, names map[string]struct{}) {
	jobs := make(chan model.Metrics, len(names))
	for w := 1; w <= a.cfg.RateLimit; w++ {
		go a.worker(w, jobs)
	}
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
		jobs <- foundMetric
	}
	close(jobs)
}

// retryableSend sends an HTTP request with retries.
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
		// Reset the request body for retries.
		req.Body.Close() //nolint:all
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
			req.Body.Close() //nolint:all
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
	}
	return nil, nil
}

// worker processes metrics and sends them to the server.
func (a *Agent) worker(id int, jobs <-chan model.Metrics) {
	for j := range jobs {
		a.logger.Debugf("worker #%v is sending %v", id, j.ID)
		metricUpdateURL := fmt.Sprintf("http://%v/update/", a.cfg.Address)

		reqBuilder := NewRequestBuilder().SetURL(metricUpdateURL).AddJSONBody(j).Sign(a.cfg.Key).Compress().SetMethod(http.MethodPost)
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
}
