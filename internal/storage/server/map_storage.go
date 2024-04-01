package storage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (s *MapStorage) StoreMetrics(path string) error {

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close() //nolint:all
	//jsonData, err := json.MarshalIndent(s.metrics, "", "   ")
	s.mu.Lock()
	defer s.mu.Unlock()
	jsonData, err := json.Marshal(s.metrics)
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	return err
}

func (s *MapStorage) RestoreMetrics(path string) error {
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
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Unmarshal(data, &s.metrics)
}
