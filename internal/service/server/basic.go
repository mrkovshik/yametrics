package server

import (
	"bytes"
	"context"
	"fmt"

	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/templates"
	"go.uber.org/zap"
)

type (
	metricService struct {
		storage service.Storage
		config  *config.ServerConfig
		logger  *zap.SugaredLogger
	}
)

func NewMetricService(storage service.Storage, config *config.ServerConfig, logger *zap.SugaredLogger) api.Service {
	return &metricService{
		storage: storage,
		config:  config,
		logger:  logger,
	}
}

func (s *metricService) UpdateMetrics(ctx context.Context, batch []model.Metrics) error {
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
func (s *metricService) GetMetric(ctx context.Context, metricModel model.Metrics) (model.Metrics, error) {
	metric, err := s.storage.GetMetricByModel(ctx, metricModel)
	if err != nil {
		errMsg := fmt.Errorf("GetMetricByModel: %s", err.Error())
		s.logger.Error(errMsg)
		return model.Metrics{}, errMsg
	}
	return metric, nil
}
func (s *metricService) GetAllMetrics(ctx context.Context) (string, error) {
	var tpl bytes.Buffer
	t, err := templates.ParseTemplates()
	if err != nil {
		return "", err
	}
	metricMap, err := s.GetAllMetrics(ctx)
	if err != nil {
		return "", err
	}
	if err := t.ExecuteTemplate(&tpl, "list_metrics", metricMap); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (s *metricService) Ping(ctx context.Context) error {

	return s.storage.Ping(ctx)
}
