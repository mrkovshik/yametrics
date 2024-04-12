package storage

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/templates"
)

type mapStorage struct {
	Mu      sync.Mutex
	Metrics map[string]model.Metrics
}

func NewMapStorage() service.Storage {
	s := make(map[string]model.Metrics)
	return &mapStorage{
		sync.Mutex{},
		s,
	}
}

func (s *mapStorage) UpdateMetricValue(_ context.Context, newMetrics model.Metrics) error {
	key := fmt.Sprintf("%v:%v", newMetrics.MType, newMetrics.ID)
	s.Mu.Lock()
	defer s.Mu.Unlock()
	found, ok := s.Metrics[key]
	if ok && (newMetrics.MType == model.MetricTypeCounter) {
		newDelta := *s.Metrics[key].Delta + *newMetrics.Delta
		found.Delta = &newDelta
		s.Metrics[key] = found
		return nil
	}
	s.Metrics[key] = newMetrics
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
	s.Mu.Lock()
	defer s.Mu.Unlock()
	res, ok := s.Metrics[key]
	if !ok {
		return model.Metrics{}, fmt.Errorf("%v not found", key)
	}

	return res, nil
}

func (s *mapStorage) GetAllMetrics(_ context.Context) (string, error) {
	var tpl bytes.Buffer
	t, err := templates.ParseTemplates()
	if err != nil {
		return "", err
	}
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if err := t.ExecuteTemplate(&tpl, "list_metrics", s.Metrics); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (s *mapStorage) StoreMetrics(_ context.Context, path string) error {

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:all
	s.Mu.Lock()
	defer s.Mu.Unlock()
	jsonData, err := json.Marshal(s.Metrics)
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	return err
}

func (s *mapStorage) RestoreMetrics(_ context.Context, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
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
	s.Mu.Lock()
	defer s.Mu.Unlock()
	return json.Unmarshal(data, &s.Metrics)
}
