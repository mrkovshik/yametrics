package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/util/retriable"
)

type mapStorage struct {
	mu      sync.RWMutex
	metrics map[string]model.Metrics
}

func NewMapStorage() service.Storage {
	s := make(map[string]model.Metrics)
	return &mapStorage{
		sync.RWMutex{},
		s,
	}
}

func (s *mapStorage) UpdateMetricValue(_ context.Context, newMetrics model.Metrics) error {
	key := fmt.Sprintf("%v:%v", newMetrics.MType, newMetrics.ID)
	s.mu.Lock()
	defer s.mu.Unlock()
	found, ok := s.metrics[key]
	if ok && (newMetrics.MType == model.MetricTypeCounter) {
		newDelta := *s.metrics[key].Delta + *newMetrics.Delta
		found.Delta = &newDelta
		s.metrics[key] = found
		return nil
	}
	s.metrics[key] = newMetrics
	return nil
}
func (s *mapStorage) UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error {
	for _, metric := range newMetrics {
		if err := s.UpdateMetricValue(ctx, metric); err != nil {
			return err
		}
	}
	return nil
}
func (s *mapStorage) GetMetricByModel(_ context.Context, newMetrics model.Metrics) (model.Metrics, error) {
	key := fmt.Sprintf("%v:%v", newMetrics.MType, newMetrics.ID)
	s.mu.RLock()
	defer s.mu.RUnlock()
	res, ok := s.metrics[key]
	if !ok {
		return model.Metrics{}, fmt.Errorf("%v not found", key)
	}

	return res, nil
}

func (s *mapStorage) GetAllMetrics(_ context.Context) (map[string]model.Metrics, error) {
	return s.metrics, nil
}

func (s *mapStorage) StoreMetrics(_ context.Context, path string) error {
	file, err := retriable.OpenRetryable(func() (*os.File, error) {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	})
	if err != nil {
		return err
	}
	defer file.Close() //nolint:all
	s.mu.RLock()
	defer s.mu.RUnlock()
	jsonData, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	return err
}

func (s *mapStorage) RestoreMetrics(_ context.Context, path string) error {
	file, err := retriable.OpenRetryable(func() (*os.File, error) {
		return os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	})
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Unmarshal(data, &s.metrics)
}

func (s *mapStorage) Ping(_ context.Context) error {
	return nil
}
