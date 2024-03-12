package storage

import (
	"errors"
	"fmt"
	"log"
	"strconv"
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

func (m *AgentMapStorage) SaveMetric(name string, value string) {
	m.Map.Store(name, value)
}

func (m *AgentMapStorage) LoadMetric(name string) string {
	value, ok := m.Map.Load(name)
	if !ok {
		return "0"
	}
	return value.(string)
}

func (m *AgentMapStorage) UpdateCounter() error {
	var (
		intValue int
		err      error
	)

	value, ok := m.Map.Load("PollCount")
	if !ok {
		intValue = 0
	} else {
		stringValue, ok := value.(string)
		if !ok {
			return errors.New("invalid server data")
		}
		intValue, err = strconv.Atoi(stringValue)
		if err != nil {
			log.Fatal(err)
		}

	}
	intValue++
	m.Map.Store("PollCount", fmt.Sprint(intValue))
	return nil
}
