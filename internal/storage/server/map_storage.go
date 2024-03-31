package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/templates"
)

type MapStorage struct {
	mu      sync.Mutex
	metrics map[string]model.Metrics
}

func NewMapStorage() *MapStorage {
	s := make(map[string]model.Metrics)
	return &MapStorage{
		sync.Mutex{},
		s,
	}
}

func (s *MapStorage) UpdateMetricValue(newMetrics model.Metrics) {
	key := fmt.Sprintf("%v:%v", newMetrics.MType, newMetrics.ID)
	s.mu.Lock()
	defer s.mu.Unlock()
	found, ok := s.metrics[key]
	if ok && (newMetrics.MType == model.MetricTypeCounter) {
		newDelta := *s.metrics[key].Delta + *newMetrics.Delta
		found.Delta = &newDelta
		s.metrics[key] = found
		return
	}
	s.metrics[key] = newMetrics

}

func (s *MapStorage) GetMetricValue(newMetrics model.Metrics) (model.Metrics, error) {
	key := fmt.Sprintf("%v:%v", newMetrics.MType, newMetrics.ID)
	s.mu.Lock()
	defer s.mu.Unlock()
	res, ok := s.metrics[key]
	if !ok {
		return model.Metrics{}, errors.New("not found")
	}

	return res, nil
}

func (s *MapStorage) GetAllMetrics() (string, error) {
	var tpl bytes.Buffer
	t, err := templates.ParseTemplates()
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := t.ExecuteTemplate(&tpl, "list_metrics", s.metrics); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (s *MapStorage) DumpMetrics(path string) error {
	jsonData, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonData, 0666)
}
