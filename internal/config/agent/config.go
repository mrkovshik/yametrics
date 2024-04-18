package config

import (
	"errors"
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/mrkovshik/yametrics/internal/util"
)

type AgentConfig struct {
	Key            string `env:"KEY"`
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

type AgentConfigBuilder struct {
	Config AgentConfig
}

func (c *AgentConfigBuilder) WithKey(key string) *AgentConfigBuilder {
	c.Config.Key = key
	return c
}

func (c *AgentConfigBuilder) WithAddress(host string) *AgentConfigBuilder {
	c.Config.Address = host
	return c
}

func (c *AgentConfigBuilder) WithReportInterval(reportInterval int) *AgentConfigBuilder {
	c.Config.ReportInterval = reportInterval
	return c
}

func (c *AgentConfigBuilder) WithPollInterval(pollInterval int) *AgentConfigBuilder {
	c.Config.PollInterval = pollInterval
	return c
}

func (c *AgentConfigBuilder) FromFlags() *AgentConfigBuilder {
	addr := flag.String("a", "localhost:8080", "server host and port")
	pollInterval := flag.Int("p", 2, "metrics polling interval")
	reportInterval := flag.Int("r", 10, "metrics sending to server interval")
	key := flag.String("k", "", "secret auth key")
	flag.Parse()

	if c.Config.Key == "" {
		c.WithKey(*key)
	}
	if c.Config.Address == "" {
		c.WithAddress(*addr)
	}
	if c.Config.PollInterval == 0 {
		c.WithPollInterval(*pollInterval)
	}
	if c.Config.ReportInterval == 0 {
		c.WithReportInterval(*reportInterval)
	}
	return c
}

func (c *AgentConfigBuilder) FromEnv() *AgentConfigBuilder {
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
	return c
}

func GetConfigs() (AgentConfig, error) {
	var c AgentConfigBuilder
	c.FromEnv().FromFlags()
	if !util.ValidateAddress(c.Config.Address) {
		return AgentConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}
