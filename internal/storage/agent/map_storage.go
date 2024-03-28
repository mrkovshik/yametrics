package storage

import (
	"errors"
	"sync"
)

type AgentMapStorage struct {
	Map sync.Map
}

func NewAgentMapStorage() *AgentMapStorage {
	return &AgentMapStorage{
		Map: sync.Map{},
	}
}

func (m *AgentMapStorage) SaveMetric(name string, value float64) {
	m.Map.Store(name, value)
}

func (m *AgentMapStorage) LoadMetric(name string) (float64, error) {
	value, ok := m.Map.Load(name)
	if !ok {
		return 0, nil
	}
	floatVal, ok := value.(float64)
	if !ok {
		return 0, errors.New("invalid server data")
	}
	return floatVal, nil
}

func (m *AgentMapStorage) LoadCounter() (int64, error) {
	value, ok := m.Map.Load("PollCount")
	if !ok {
		return 0, nil
	}
	intValue, ok := value.(int64)
	if !ok {
		return 0, errors.New("invalid server data")
	}
	return intValue, nil
}

func (m *AgentMapStorage) UpdateCounter() error {
	intValue, err := m.LoadCounter()
	if err != nil {
		return err
	}
	intValue++
	m.Map.Store("PollCount", intValue)
	return nil
}
