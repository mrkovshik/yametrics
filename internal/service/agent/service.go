// Package service provides utilities for building and sending HTTP requests with metrics to server.
package service

import (
	"context"
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

type sender interface {
	Send(id int, jobs <-chan model.Metrics)
}

// Agent represents a metric collection agent that polls and sends metrics.
type Agent struct {
	client  sender               //Client for sending metrics to the server
	source  metrics.MetricSource // Source of the metrics
	logger  *zap.SugaredLogger   // Logger for logging messages
	cfg     *config.AgentConfig  // Configuration for the agent
	storage storage              // Storage for metrics
}

// NewAgent initializes a new Agent.
func NewAgent(source metrics.MetricSource, cfg *config.AgentConfig, strg storage, logger *zap.SugaredLogger, client sender) *Agent {
	return &Agent{
		client:  client,
		source:  source,
		logger:  logger,
		cfg:     cfg,
		storage: strg,
	}
}

// SendMetrics sends metrics at intervals specified by the channel.
func (a *Agent) SendMetrics(ctx context.Context, ch <-chan time.Time, done chan struct{}) {
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
		a.logger.Debug("Starting to send metrics")
		a.sendMetricsByPool(ctx, metricNamesMap)
		a.logger.Debug("Metrics sent.\n")
	}
	done <- struct{}{}
}

// PollMetrics polls metrics at intervals specified by the channel.
func (a *Agent) PollMetrics(ch <-chan time.Time, done chan struct{}) {
	for range ch {
		a.logger.Debug("Starting to update metrics")
		if err := a.source.PollMemStats(a.storage); err != nil {
			a.logger.Error("PollMemStats", err)
			return
		}
		a.logger.Debug("Metrics updated.\n")
	}
	done <- struct{}{}
}

// PollUtilMetrics polls utilization metrics at intervals specified by the channel.
func (a *Agent) PollUtilMetrics(ch <-chan time.Time, done chan struct{}) {
	for range ch {
		if err := a.source.PollVirtMemStats(a.storage); err != nil {
			a.logger.Error("PollVirtMemStats", err)
			return
		}
		a.logger.Debug("Metrics updated.\n")
	}
	done <- struct{}{}
}

// sendMetricsByPool sends metrics using a pool of workers.
func (a *Agent) sendMetricsByPool(ctx context.Context, names map[string]struct{}) {
	jobs := make(chan model.Metrics, len(names))
	for w := 1; w <= a.cfg.RateLimit; w++ {
		go a.client.Send(w, jobs)
	}
	for name := range names {
		currentMetric := model.Metrics{
			ID: name,
		}
		switch name {
		case "PollCount":
			currentMetric.MType = model.MetricTypeCounter
		default:
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
