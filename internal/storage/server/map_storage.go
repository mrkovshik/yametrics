package storage

import (
	"bytes"
	"errors"
	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/templates"
)

type MapStorage map[string]model.Metrics

func NewMapStorage() MapStorage {
	s := make(map[string]model.Metrics)
	return s
}

func (s MapStorage) UpdateMetricValue(newMetrics model.Metrics) {
	s[newMetrics.ID] = newMetrics
}

func (s MapStorage) GetMetricValue(newMetrics model.Metrics) (model.Metrics, error) {

	res, ok := s[newMetrics.ID]
	if !ok {
		return model.Metrics{}, errors.New("not found")
	}
	return res, nil

}

func (s MapStorage) GetAllMetrics() (string, error) {
	var tpl bytes.Buffer
	t, err := templates.ParseTemplates()
	if err != nil {
		return "", err
	}
	if err := t.ExecuteTemplate(&tpl, "list_metrics", s); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
