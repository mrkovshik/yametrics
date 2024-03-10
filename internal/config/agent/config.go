package config

import (
	"errors"
	"flag"
	"github.com/mrkovshik/yametrics/internal/utl"
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func (c *AgentConfig) GetConfigs() error {
	addr := flag.String("a", "localhost:8080", "server host and port")
	pollInterval := flag.Int("p", 2, "metrics polling interval")
	reportInterval := flag.Int("r", 10, "metrics sending to server interval")
	flag.Parse()
	c.Address = *addr
	if !utl.ValidateAddress(c.Address) {
		return errors.New("need address in a form host:port")
	}
	c.ReportInterval = *reportInterval
	c.PollInterval = *pollInterval
	return nil
}
