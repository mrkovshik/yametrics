package config

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/mrkovshik/yametrics/internal/utl"
	"log"
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
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
	if c.Address == "" {
		c.Address = *addr
	}
	if !utl.ValidateAddress(c.Address) {
		return errors.New("need address in a form host:port")
	}
	if c.ReportInterval == 0 {
		c.ReportInterval = *reportInterval
	}
	if c.PollInterval == 0 {
		c.PollInterval = *pollInterval
	}
	return nil
}
