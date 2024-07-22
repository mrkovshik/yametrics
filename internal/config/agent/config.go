// Package config provides configuration handling for the agent, allowing
// configurations to be set via environment variables or command-line flags.
package config

import (
	"errors"
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/eschao/config"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mrkovshik/yametrics/internal/util"
)

const (
	defaultKey            = ""
	defaultAddress        = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultRateLimit      = 1
	defaultCryptoKey      = "./public_key.pem"
)

var k = koanf.New(".")

// AgentConfig holds the configuration settings for the agent.
type AgentConfig struct {
	Key            string `env:"KEY" json:"key"`
	Address        string `env:"ADDRESS" json:"address"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	RateLimit      int    `env:"RATE_LIMIT" json:"rate_limit"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFilePath string `env:"CONFIG" json:"config_file_path"`
}

// AgentJSON is a structure fot mapping config from JSON file.
type AgentJSON struct {
	Key            string `json:"key"`
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
	RateLimit      int    `json:"rate_limit"`
	CryptoKey      string `json:"crypto_key"`
	ConfigFilePath string `json:"config_file_path"`
}

// AgentConfigBuilder is a builder for constructing an AgentConfig instance.
type AgentConfigBuilder struct {
	Config AgentConfig
}

// WithKey sets the key in the AgentConfig.
func (c *AgentConfigBuilder) WithKey(key string) *AgentConfigBuilder {
	c.Config.Key = key
	return c
}

// WithAddress sets the address in the AgentConfig.
func (c *AgentConfigBuilder) WithAddress(address string) *AgentConfigBuilder {
	c.Config.Address = address
	return c
}

// WithCryptoKey sets the crypto key flag in the ServerConfig.
func (c *AgentConfigBuilder) WithCryptoKey(path string) *AgentConfigBuilder {
	c.Config.CryptoKey = path
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
	c.Config.RateLimit = rateLimit
	return c
}

// WithConfigFile sets the path to JSON configuration file
func (c *AgentConfigBuilder) WithConfigFile(configFilePath string) *AgentConfigBuilder {
	c.Config.ConfigFilePath = configFilePath
	return c
}

// FromFlags populates the AgentConfig from command-line flags.
func (c *AgentConfigBuilder) FromFlags() *AgentConfigBuilder {
	addr := flag.String("a", defaultAddress, "server host and port")
	pollInterval := flag.Int("p", defaultPollInterval, "metrics polling interval")
	reportInterval := flag.Int("r", defaultReportInterval, "metrics sending to server interval")
	key := flag.String("k", defaultKey, "secret auth key")
	rateLimit := flag.Int("l", defaultRateLimit, "agent rate limit")
	cryptoKey := flag.String("crypto-key", defaultCryptoKey, "path to the file with public key")
	configFilePath := flag.String("c", "", "path to config file")
	configFilePathAlias := flag.String("config", "", "path to config file (shorthand)")
	flag.Parse()

	if *configFilePath != "" && *configFilePathAlias != "" {
		log.Fatalf("usage of both shorthand and full flag (-c and --config)")
	}

	if c.Config.ConfigFilePath == "" {
		if *configFilePath != "" {
			c.WithConfigFile(*configFilePath)
		}
		if *configFilePathAlias != "" {
			c.WithConfigFile(*configFilePathAlias)
		}
	}

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
	if c.Config.RateLimit == 0 {
		c.WithRateLimit(*rateLimit)
	}
	if c.Config.CryptoKey == "" {
		c.WithCryptoKey(*cryptoKey)
	}
	return c
}

// FromFile populates the AgentConfig from config JSON file.
func (c *AgentConfigBuilder) FromFile() *AgentConfigBuilder {
	if c.Config.ConfigFilePath != "" {
		if err := k.Load(file.Provider(c.Config.ConfigFilePath), json.Parser()); err != nil {
			log.Fatalf("error loading config: %v", err)
		}
	}
	JSONConfig := AgentJSON{}
	if err := config.ParseConfigFile(&JSONConfig, c.Config.ConfigFilePath); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}
	if JSONConfig.Key != "" && c.Config.Key == defaultKey {
		c.WithKey(JSONConfig.Key)
	}
	if JSONConfig.Address != "" && c.Config.Address == defaultAddress {
		c.WithAddress(JSONConfig.Address)
	}
	if JSONConfig.ReportInterval != "" && c.Config.ReportInterval == defaultReportInterval {
		intReportInterval, err := util.CutSeconds(JSONConfig.ReportInterval)
		if err != nil {
			log.Fatal(err)
		}
		c.WithReportInterval(intReportInterval)
	}
	if JSONConfig.PollInterval != "" && c.Config.PollInterval == defaultPollInterval {
		intPollInterval, err := util.CutSeconds(JSONConfig.PollInterval)
		if err != nil {
			log.Fatal(err)
		}
		c.WithPollInterval(intPollInterval)
	}
	if JSONConfig.RateLimit != 0 && c.Config.RateLimit == defaultRateLimit {
		c.WithRateLimit(JSONConfig.RateLimit)
	}
	if JSONConfig.CryptoKey != "" && c.Config.CryptoKey == defaultCryptoKey {
		c.WithCryptoKey(JSONConfig.CryptoKey)
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
// configurations from environment variables, command-line flags and JSON file.
// It validates the address and rate limit to ensure they are properly set.
func GetConfigs() (AgentConfig, error) {
	var c AgentConfigBuilder
	c.FromEnv().FromFlags().FromFile()
	if !util.ValidateAddress(c.Config.Address) {
		return AgentConfig{}, errors.New("need address in a form host:port")
	}
	if c.Config.RateLimit == 0 {
		return AgentConfig{}, errors.New("rate limit must be larger than 0")
	}
	return c.Config, nil
}
