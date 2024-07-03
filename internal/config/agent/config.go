// Package config provides configuration handling for the agent, allowing
// configurations to be set via environment variables or command-line flags.
package config

import (
	"errors"
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	service "github.com/mrkovshik/yametrics/internal/service/agent"
	"github.com/mrkovshik/yametrics/internal/util"
)

// AgentConfig holds the configuration settings for the agent.
type AgentConfig struct {
	ServiceConfig  service.Config
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

// AgentConfigBuilder is a builder for constructing an AgentConfig instance.
type AgentConfigBuilder struct {
	Config AgentConfig
}

// WithKey sets the key in the AgentConfig.
func (c *AgentConfigBuilder) WithKey(key string) *AgentConfigBuilder {
	c.Config.ServiceConfig.Key = key
	return c
}

// WithAddress sets the address in the AgentConfig.
func (c *AgentConfigBuilder) WithAddress(address string) *AgentConfigBuilder {
	c.Config.ServiceConfig.Address = address
	return c
}

// WithReportInterval sets the report interval in the AgentConfig.
func (c *AgentConfigBuilder) WithReportInterval(reportInterval int) *AgentConfigBuilder {
	c.Config.ReportInterval = reportInterval
	return c
}

// WithPollInterval sets the poll interval in the AgentConfig.
func (c *AgentConfigBuilder) WithPollInterval(pollInterval int) *AgentConfigBuilder {
	c.Config.PollInterval = pollInterval
	return c
}

// WithRateLimit sets the rate limit in the AgentConfig.
func (c *AgentConfigBuilder) WithRateLimit(rateLimit int) *AgentConfigBuilder {
	c.Config.ServiceConfig.RateLimit = rateLimit
	return c
}

// FromFlags populates the AgentConfig from command-line flags.
func (c *AgentConfigBuilder) FromFlags() *AgentConfigBuilder {
	addr := flag.String("a", "localhost:8080", "server host and port")
	pollInterval := flag.Int("p", 2, "metrics polling interval")
	reportInterval := flag.Int("r", 10, "metrics sending to server interval")
	key := flag.String("k", "", "secret auth key")
	rateLimit := flag.Int("l", 1, "agent rate limit")
	flag.Parse()

	if c.Config.ServiceConfig.Key == "" {
		c.WithKey(*key)
	}
	if c.Config.ServiceConfig.Address == "" {
		c.WithAddress(*addr)
	}
	if c.Config.PollInterval == 0 {
		c.WithPollInterval(*pollInterval)
	}
	if c.Config.ReportInterval == 0 {
		c.WithReportInterval(*reportInterval)
	}
	if c.Config.ServiceConfig.RateLimit == 0 {
		c.WithRateLimit(*rateLimit)
	}
	return c
}

// FromEnv populates the AgentConfig from environment variables.
func (c *AgentConfigBuilder) FromEnv() *AgentConfigBuilder {
	if err := env.Parse(&c.Config); err != nil {
		log.Fatal(err)
	}
	return c
}

// GetConfigs returns the fully constructed AgentConfig by combining
// configurations from environment variables and command-line flags.
// It validates the address and rate limit to ensure they are properly set.
func GetConfigs() (AgentConfig, error) {
	var c AgentConfigBuilder
	c.FromEnv().FromFlags()
	if !util.ValidateAddress(c.Config.ServiceConfig.Address) {
		return AgentConfig{}, errors.New("need address in a form host:port")
	}
	if c.Config.ServiceConfig.RateLimit == 0 {
		return AgentConfig{}, errors.New("rate limit must be larger than 0")
	}
	return c.Config, nil
}
