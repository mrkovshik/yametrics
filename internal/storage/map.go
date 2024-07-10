// Package storage provides implementations of the service.Storage interface for metrics storage.
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
	"github.com/mrkovshik/yametrics/internal/util/retriable"
)

// InMemoryStorage implements the service.Storage interface using an in-memory map for storing metrics.
type InMemoryStorage struct {
	mu      sync.RWMutex             // Mutex for thread-safe access to metrics map
	metrics map[string]model.Metrics // Map to store metrics
}

// NewInMemoryStorage creates a new instance of InMemoryStorage.
// Returns:
// - a pointer to the new InMemoryStorage instance.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		metrics: make(map[string]model.Metrics),
	}
}

// UpdateMetricValue updates or inserts a metric into the metrics map.
// Parameters:
// - ctx: the context to control the update operation.
// - newMetrics: the Metrics model containing the metric data to be updated or inserted.
// Returns:
// - an error if the update operation fails.
func (s *InMemoryStorage) UpdateMetricValue(_ context.Context, newMetrics model.Metrics) error {
	key := newMetrics.MType + ":" + newMetrics.ID
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

// UpdateMetrics updates multiple metrics in the metrics map.
// Parameters:
// - ctx: the context to control the update operation.
// - newMetrics: a slice of Metrics models containing the metric data to be updated or inserted.
// Returns:
// - an error if the update operation fails.
func (s *InMemoryStorage) UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error {
	for _, metric := range newMetrics {
		if err := s.UpdateMetricValue(ctx, metric); err != nil {
			return err
		}
	}
	return nil
}

// GetMetricByModel retrieves a metric from the metrics map based on the provided model.
// Parameters:
// - ctx: the context to control the retrieval operation.
// - newMetrics: the Metrics model specifying the metric to be retrieved.
// Returns:
// - the retrieved Metrics model.
// - an error if the retrieval operation fails.
func (s *InMemoryStorage) GetMetricByModel(_ context.Context, newMetrics model.Metrics) (model.Metrics, error) {
	key := newMetrics.MType + ":" + newMetrics.ID
	s.mu.RLock()
	defer s.mu.RUnlock()
	res, ok := s.metrics[key]
	if !ok {
		return model.Metrics{}, fmt.Errorf("%v not found", key)
	}
	return res, nil
}

// GetAllMetrics retrieves all metrics from the metrics map.
// Parameters:
// - ctx: the context to control the retrieval operation.
// Returns:
// - a map of metric names to Metrics models representing all stored metrics.
// - an error if the retrieval operation fails.
func (s *InMemoryStorage) GetAllMetrics(_ context.Context) (map[string]model.Metrics, error) {
	//newMap := make(map[string]model.Metrics, len(s.metrics))
	//for _,v:=range s.metrics{
	//	newMap[id]
	//}
	return s.metrics, nil
}

// StoreMetrics stores all metrics from the metrics map into a JSON file at the specified path.
// Parameters:
// - ctx: the context to control the store operation.
// - path: the file path where metrics should be stored.
// Returns:
// - an error if the store operation fails.
func (s *InMemoryStorage) StoreMetrics(_ context.Context, path string) error {
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

// RestoreMetrics restores metrics from a JSON file at the specified path into the metrics map.
// Parameters:
// - ctx: the context to control the restore operation.
// - path: the file path from where metrics should be restored.
// Returns:
// - an error if the restore operation fails.
func (s *InMemoryStorage) RestoreMetrics(_ context.Context, path string) error {
	file, err := retriable.OpenRetryable(func() (*os.File, error) {
		return os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	})
	if err != nil {
		return err
	}
	defer file.Close() //nolint:all
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

// Ping checks the availability of the InMemoryStorage.
// Parameters:
// - ctx: the context to control the ping operation.
// Returns:
// - an error if the storage is not available.
func (s *InMemoryStorage) Ping(_ context.Context) error {
	return nil
}
