// Package server provides the implementation of the metric service
// which includes methods to update, retrieve, and store metrics.
package server

import (
	"bytes"
	"context"
	"fmt"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/templates"
	"go.uber.org/zap"
)

type storage interface {
	UpdateMetricValue(ctx context.Context, newMetrics model.Metrics) error

	UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error

	GetMetricByModel(ctx context.Context, newMetrics model.Metrics) (model.Metrics, error)

	GetAllMetrics(ctx context.Context) (map[string]model.Metrics, error)

	StoreMetrics(ctx context.Context, path string) error

	RestoreMetrics(ctx context.Context, path string) error

	Ping(ctx context.Context) error
}

// MetricService represents the service for managing metrics.
type MetricService struct {
	storage storage
	config  *config.ServerConfig
	logger  *zap.SugaredLogger
}

// NewMetricService creates a new instance of MetricService.
//
// storage: an implementation of the Storage interface for managing metric data.
// config: the server configuration settings.
// logger: a logger for logging messages.
func NewMetricService(storage storage, config *config.ServerConfig, logger *zap.SugaredLogger) *MetricService {
	return &MetricService{
		storage: storage,
		config:  config,
		logger:  logger,
	}
}

// UpdateMetrics updates the metrics in the storage. If SyncStoreEnable is true in the config,
// it also stores the metrics to the file specified in StoreFilePath.
//
// ctx: the context for managing request-scoped values and cancelation.
// batch: a slice of metrics to be updated.
//
// Returns an error if the update or store operation fails.
func (s *MetricService) UpdateMetrics(ctx context.Context, batch []model.Metrics) error {
	if err := s.storage.UpdateMetrics(ctx, batch); err != nil {
		errMsg := fmt.Errorf("UpdateMetrics: %s", err.Error())
		s.logger.Error(errMsg)
		return errMsg
	}
	if s.config.SyncStoreEnable {
		if err := s.storage.StoreMetrics(ctx, s.config.StoreFilePath); err != nil {
			errMsg := fmt.Errorf("StoreMetrics: %s", err.Error())
			s.logger.Error(errMsg)
			return errMsg
		}
	}
	return nil
}

// GetMetric retrieves a specific metric from the storage based on the provided metric model.
//
// ctx: the context for managing request-scoped values and cancelation.
// metricModel: the model of the metric to be retrieved.
//
// Returns the retrieved metric and an error if the retrieval fails.
func (s *MetricService) GetMetric(ctx context.Context, metricModel model.Metrics) (model.Metrics, error) {
	metric, err := s.storage.GetMetricByModel(ctx, metricModel)
	if err != nil {
		errMsg := fmt.Errorf("GetMetricByModel: %s", err.Error())
		s.logger.Error(errMsg)
		return model.Metrics{}, errMsg
	}
	return metric, nil
}

// GetAllMetrics retrieves all metrics from the storage and returns them as a formatted string.
//
// ctx: the context for managing request-scoped values and cancelation.
//
// Returns a formatted string of all metrics and an error if the retrieval or template execution fails.
func (s *MetricService) GetAllMetrics(ctx context.Context) (string, error) {
	var tpl bytes.Buffer
	t, err := templates.ParseTemplates()
	if err != nil {
		return "", err
	}
	metricMap, err := s.storage.GetAllMetrics(ctx)
	if err != nil {
		return "", err
	}
	if err := t.ExecuteTemplate(&tpl, "list_metrics", metricMap); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (s *MetricService) StoreMetrics(ctx context.Context) error {
	return s.storage.StoreMetrics(ctx, s.config.StoreFilePath)
}

func (s *MetricService) RestoreMetrics(ctx context.Context) error {
	return s.storage.RestoreMetrics(ctx, s.config.StoreFilePath)
}

// Ping checks the connectivity to the storage.
//
// ctx: the context for managing request-scoped values and cancelation.
//
// Returns an error if the ping operation fails.
func (s *MetricService) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}
