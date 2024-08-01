// Package config provides configuration handling for the agent, allowing
// configurations to be set via environment variables or command-line flags.
package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/eschao/config"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mrkovshik/yametrics/internal/config/flags"
	"github.com/mrkovshik/yametrics/internal/util"
)

const (
	defaultKey            = ""
	defaultConfigFilePath = ""
	defaultAddress        = "localhost:8080"
	defaultPollInterval   = 2
	defaultReportInterval = 10
	defaultRateLimit      = 1
	defaultCryptoKey      = "./public_key.pem"
)

var k = koanf.New(".")

// AgentConfig holds the configuration settings for the agent.
// TODO: divide config for client and metric
type AgentConfig struct {
	Key                  string `env:"KEY" json:"key"`
	KeyIsSet             bool   `json:"-"`
	Address              string `env:"ADDRESS" json:"address"`
	AddressIsSet         bool   `json:"-"`
	ReportInterval       int    `env:"REPORT_INTERVAL"`
	ReportIntervalString string `json:"report_interval"`
	ReportIntervalIsSet  bool   `json:"-"`
	PollInterval         int    `env:"POLL_INTERVAL"`
	PollIntervalString   string `json:"poll_interval"`
	PollIntervalIsSet    bool   `json:"-"`
	RateLimit            int    `env:"RATE_LIMIT" json:"rate_limit"`
	RateLimitIsSet       bool   `json:"-"`
	CryptoKey            string `env:"CRYPTO_KEY" json:"crypto_key"`
	CryptoKeyIsSet       bool   `json:"-"`
	ConfigFilePath       string `env:"CONFIG" json:"config_file_path"`
	ConfigFilePathIsSet  bool   `json:"-"`
}

// AgentConfigBuilder is a builder for constructing an AgentConfig instance.
type AgentConfigBuilder struct {
	Config AgentConfig
}

func (c *AgentConfig) SetDefaults() {
	c.Key = defaultKey
	c.Address = defaultAddress
	c.CryptoKey = defaultCryptoKey
	c.RateLimit = defaultRateLimit
	c.ReportInterval = defaultReportInterval
	c.PollInterval = defaultPollInterval
	c.ConfigFilePath = defaultConfigFilePath
}

// WithKey sets the key in the AgentConfig.
func (c *AgentConfigBuilder) WithKey(key string) *AgentConfigBuilder {
	c.Config.Key = key
	c.Config.KeyIsSet = true
	return c
}

// WithAddress sets the address in the AgentConfig.
func (c *AgentConfigBuilder) WithAddress(address string) *AgentConfigBuilder {
	c.Config.Address = address
	c.Config.AddressIsSet = true
	return c
}

// WithCryptoKey sets the crypto key flag in the ServerConfig.
func (c *AgentConfigBuilder) WithCryptoKey(path string) *AgentConfigBuilder {
	c.Config.CryptoKey = path
	c.Config.CryptoKeyIsSet = true
	return c
}

// WithReportInterval sets the report interval in the AgentConfig.
func (c *AgentConfigBuilder) WithReportInterval(reportInterval int) *AgentConfigBuilder {
	c.Config.ReportInterval = reportInterval
	c.Config.ReportIntervalIsSet = true
	return c
}

// WithPollInterval sets the poll interval in the AgentConfig.
func (c *AgentConfigBuilder) WithPollInterval(pollInterval int) *AgentConfigBuilder {
	c.Config.PollInterval = pollInterval
	c.Config.PollIntervalIsSet = true
	return c
}

// WithRateLimit sets the rate limit in the AgentConfig.
func (c *AgentConfigBuilder) WithRateLimit(rateLimit int) *AgentConfigBuilder {
	c.Config.RateLimit = rateLimit
	c.Config.RateLimitIsSet = true
	return c
}

// WithConfigFile sets the path to JSON configuration file
func (c *AgentConfigBuilder) WithConfigFile(configFilePath string) *AgentConfigBuilder {
	c.Config.ConfigFilePath = configFilePath
	c.Config.ConfigFilePathIsSet = true
	return c
}

// FromFlags populates the AgentConfig from command-line flags.
func (c *AgentConfigBuilder) FromFlags() *AgentConfigBuilder {
	addr := flags.CustomString{}
	flag.Var(&addr, "a", "server host and port")

	pollInterval := flags.CustomInt{}
	flag.Var(&pollInterval, "p", "metrics polling interval")

	reportInterval := flags.CustomInt{}
	flag.Var(&reportInterval, "r", "metrics sending to server interval")

	key := flags.CustomString{}
	flag.Var(&key, "k", "secret auth key")

	rateLimit := flags.CustomInt{}
	flag.Var(&rateLimit, "l", "number of agent workers")

	cryptoKey := flags.CustomString{}
	flag.Var(&cryptoKey, "crypto-key", "path to the file with public key")

	configFilePath := flags.CustomString{}
	flag.Var(&configFilePath, "c", "path to config file (shorthand)")

	configFilePathAlias := flags.CustomString{}
	flag.Var(&configFilePathAlias, "config", "path to config file")

	flag.Parse()

	//Verifying if the flags were set properly
	if configFilePath.IsSet && configFilePathAlias.IsSet {
		log.Fatalf("usage of both shorthand and full flag (-c and --config)")
	}

	if !c.Config.ConfigFilePathIsSet {
		if configFilePath.IsSet {
			c.WithConfigFile(configFilePath.Value)
		}
		if configFilePathAlias.IsSet {
			c.WithConfigFile(configFilePathAlias.Value)
		}
	}

	if !c.Config.KeyIsSet && key.IsSet {
		c.WithKey(key.Value)
	}

	if !c.Config.CryptoKeyIsSet && cryptoKey.IsSet {
		c.WithCryptoKey(cryptoKey.Value)
	}

	if !c.Config.AddressIsSet && addr.IsSet {
		c.WithAddress(addr.Value)
	}

	if !c.Config.PollIntervalIsSet && pollInterval.IsSet {
		c.WithPollInterval(pollInterval.Value)
	}

	if !c.Config.ReportIntervalIsSet && reportInterval.IsSet {
		c.WithReportInterval(reportInterval.Value)
	}

	if !c.Config.RateLimitIsSet && rateLimit.IsSet {
		c.WithRateLimit(rateLimit.Value)
	}

	return c
}

// FromFile populates the AgentConfig from config JSON file.
func (c *AgentConfigBuilder) FromFile() *AgentConfigBuilder {
	if c.Config.ConfigFilePath == "" {
		return c
	}
	if err := k.Load(file.Provider(c.Config.ConfigFilePath), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	JSONConfig := AgentConfig{}
	JSONConfig.SetDefaults()

	if err := config.ParseConfigFile(&JSONConfig, c.Config.ConfigFilePath); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}
	if JSONConfig.Key != defaultKey && !c.Config.KeyIsSet {
		c.WithKey(JSONConfig.Key)
	}

	if JSONConfig.Address != defaultAddress && !c.Config.AddressIsSet {
		c.WithAddress(JSONConfig.Address)
	}

	if JSONConfig.CryptoKey != defaultCryptoKey && !c.Config.CryptoKeyIsSet {
		c.WithCryptoKey(JSONConfig.CryptoKey)
	}

	if JSONConfig.ReportIntervalString != "" && !c.Config.ReportIntervalIsSet {
		intReportInterval, err := util.CutSeconds(JSONConfig.ReportIntervalString)
		if err != nil {
			log.Fatal(err)
		}
		c.WithReportInterval(intReportInterval)
	}

	if JSONConfig.PollIntervalString != "" && !c.Config.PollIntervalIsSet {
		intPollInterval, err := util.CutSeconds(JSONConfig.PollIntervalString)
		if err != nil {
			log.Fatal(err)
		}
		c.WithPollInterval(intPollInterval)
	}

	if JSONConfig.RateLimit != defaultRateLimit && !c.Config.RateLimitIsSet {
		c.WithRateLimit(JSONConfig.RateLimit)
	}
	return c
}

// FromEnv populates the AgentConfig from environment variables.
func (c *AgentConfigBuilder) FromEnv() *AgentConfigBuilder {
	if err := env.Parse(&c.Config); err != nil {
		log.Fatal(err)
	}

	_, addressIsSet := os.LookupEnv("ADDRESS")
	if addressIsSet {
		c.Config.AddressIsSet = true
	}
	_, keyIsSet := os.LookupEnv("KEY")
	if keyIsSet {
		c.Config.KeyIsSet = true
	}
	_, cryptoKeyIsSet := os.LookupEnv("CRYPTO_KEY")
	if cryptoKeyIsSet {
		c.Config.CryptoKeyIsSet = true
	}
	_, configFilePathIsSet := os.LookupEnv("CONFIG")
	if configFilePathIsSet {
		c.Config.ConfigFilePathIsSet = true
	}

	_, reportIntervalIsSet := os.LookupEnv("REPORT_INTERVAL")
	if reportIntervalIsSet {
		c.Config.ReportIntervalIsSet = true
	}

	_, pollIntervalIsSet := os.LookupEnv("POLL_INTERVAL")
	if pollIntervalIsSet {
		c.Config.PollIntervalIsSet = true
	}

	_, rateLimitIsSet := os.LookupEnv("RATE_LIMIT")
	if rateLimitIsSet {
		c.Config.RateLimitIsSet = true
	}

	return c
}

// GetConfigs returns the fully constructed AgentConfig by combining
// configurations from environment variables, command-line flags and JSON file.
// It validates the address and rate limit to ensure they are properly set.
func GetConfigs() (AgentConfig, error) {
	var c AgentConfigBuilder
	c.Config.SetDefaults()

	c.FromEnv().FromFlags().FromFile()

	if !util.ValidateAddress(c.Config.Address) {
		return AgentConfig{}, errors.New("need address in a form host:port")
	}
	if c.Config.RateLimit == 0 {
		return AgentConfig{}, errors.New("rate limit must be larger than 0")
	}
	return c.Config, nil
}
