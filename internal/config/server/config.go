// Package config provides configuration handling for the server, allowing
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
	defaultStoreInterval  = 300
	defaultStoreFilePath  = "./tmp/metrics-db.json"
	defaultCryptoKey      = "./public_key.pem"
	defaultDBAddress      = ""
	defaultRestoreEnable  = true
	defaultStoreEnable    = true
)

var k = koanf.New(".")

// ServerConfig holds the configuration settings for the server.
type ServerConfig struct {
	Address             string `env:"ADDRESS" json:"address"`
	AddressIsSet        bool   `json:"-"`
	Key                 string `env:"KEY" json:"key"`
	KeyIsSet            bool   `json:"-"`
	StoreInterval       int    `env:"STORE_INTERVAL" json:"-"`
	StoreIntervalString string `json:"store_interval"`
	StoreIntervalIsSet  bool   `json:"-"`
	SyncStoreEnable     bool   `json:"-"`
	StoreFilePath       string `env:"FILE_STORAGE_PATH" json:"store_file"`
	StoreFilePathIsSet  bool   `json:"-"`
	StoreEnable         bool   `json:"-"`
	RestoreEnable       bool   `env:"RESTORE" json:"restore"`
	RestoreEnvIsSet     bool   `json:"-"`
	DBAddress           string `env:"DATABASE_DSN" json:"database_dsn"`
	DBAddressIsSet      bool   `json:"-"`
	DBEnable            bool   `json:"-"`
	CryptoKey           string `env:"CRYPTO_KEY" json:"crypto_key"`
	CryptoKeyIsSet      bool   `json:"-"`
	ConfigFilePath      string `env:"CONFIG" json:"-"`
	ConfigFilePathIsSet bool   `json:"-"`
}

// ServerConfigBuilder is a builder for constructing a ServerConfig instance.
type ServerConfigBuilder struct {
	Config ServerConfig
}

func (c *ServerConfig) SetDefaults() {
	c.Key = defaultKey
	c.Address = defaultAddress
	c.DBAddress = defaultDBAddress
	c.CryptoKey = defaultCryptoKey
	c.RestoreEnable = defaultRestoreEnable
	c.StoreInterval = defaultStoreInterval
	c.StoreFilePath = defaultStoreFilePath
	c.ConfigFilePath = defaultConfigFilePath
	c.StoreEnable = defaultStoreEnable
}

// WithKey sets the key in the ServerConfig.
func (c *ServerConfigBuilder) WithKey(key string) *ServerConfigBuilder {
	c.Config.Key = key
	c.Config.KeyIsSet = true
	return c
}

// WithAddress sets the address in the ServerConfig.
func (c *ServerConfigBuilder) WithAddress(address string) *ServerConfigBuilder {
	c.Config.Address = address
	c.Config.AddressIsSet = true
	return c
}

// WithDSN sets the database DSN in the ServerConfig.
func (c *ServerConfigBuilder) WithDSN(dsn string) *ServerConfigBuilder {
	c.Config.DBAddress = dsn
	c.Config.DBAddressIsSet = true
	if dsn == "" {
		c.Config.DBEnable = false
	} else {
		c.Config.DBEnable = true
	}
	return c
}

// WithStoreInterval sets the store interval in the ServerConfig.
func (c *ServerConfigBuilder) WithStoreInterval(interval int) *ServerConfigBuilder {
	c.Config.StoreInterval = interval
	c.Config.StoreIntervalIsSet = true
	if interval == 0 {
		c.Config.SyncStoreEnable = true
	} else {
		c.Config.SyncStoreEnable = false
	}
	return c
}

// WithStoreFilePath sets the store file path in the ServerConfig.
func (c *ServerConfigBuilder) WithStoreFilePath(path string) *ServerConfigBuilder {
	c.Config.StoreFilePath = path
	c.Config.StoreFilePathIsSet = true
	if path == "" {
		c.Config.StoreEnable = false
	} else {
		c.Config.StoreEnable = true
	}
	return c
}

// WithRestoreEnable sets the restore enable flag in the ServerConfig.
func (c *ServerConfigBuilder) WithRestoreEnable(restore bool) *ServerConfigBuilder {
	c.Config.RestoreEnable = restore
	c.Config.RestoreEnvIsSet = true
	return c
}

// WithCryptoKey sets the crypto key flag in the ServerConfig.
func (c *ServerConfigBuilder) WithCryptoKey(path string) *ServerConfigBuilder {
	c.Config.CryptoKey = path
	c.Config.CryptoKeyIsSet = true
	return c
}

// WithConfigFile sets the path to JSON configuration file
func (c *ServerConfigBuilder) WithConfigFile(configFilePath string) *ServerConfigBuilder {
	c.Config.ConfigFilePath = configFilePath
	c.Config.ConfigFilePathIsSet = true
	return c
}

// FromFlags populates the ServerConfig from command-line flags.
func (c *ServerConfigBuilder) FromFlags() *ServerConfigBuilder {
	storeInterval := flags.CustomInt{}
	flag.Var(&storeInterval, "i", "time interval between storing data to file in seconds")

	addr := flags.CustomString{}
	flag.Var(&addr, "a", "server host and port")

	storeFilePath := flags.CustomString{}
	flag.Var(&storeFilePath, "f", "path to storing data file")

	restoreEnable := flags.CustomBool{}
	flag.Var(&restoreEnable, "r", "is data restore from file enabled")

	dbAddress := flags.CustomString{}
	flag.Var(&dbAddress, "d", "db address") //

	key := flags.CustomString{}
	flag.Var(&key, "k", "secret auth key")

	cryptoKey := flags.CustomString{}
	flag.Var(&cryptoKey, "crypto-key", "path to the file with private key")

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

	if !c.Config.DBAddressIsSet && dbAddress.IsSet {
		c.WithDSN(dbAddress.Value)
	}

	if !c.Config.StoreFilePathIsSet && storeFilePath.IsSet {
		c.WithStoreFilePath(storeFilePath.Value)
	}

	if !c.Config.StoreIntervalIsSet && storeInterval.IsSet {
		c.WithStoreInterval(storeInterval.Value)
	}

	if !c.Config.RestoreEnvIsSet && restoreEnable.IsSet {
		c.WithRestoreEnable(restoreEnable.Value)
	}
	return c
}

// FromFile populates the ServerConfig from JSON .
func (c *ServerConfigBuilder) FromFile() *ServerConfigBuilder {

	if c.Config.ConfigFilePath == "" {
		return c
	}
	if err := k.Load(file.Provider(c.Config.ConfigFilePath), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	JSONConfig := ServerConfig{}
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

	if JSONConfig.DBAddress != defaultDBAddress && !c.Config.DBAddressIsSet {
		c.WithDSN(JSONConfig.DBAddress)
	}

	if JSONConfig.StoreFilePath != defaultStoreFilePath && !c.Config.StoreFilePathIsSet {
		c.WithStoreFilePath(JSONConfig.StoreFilePath)
	}

	if JSONConfig.StoreIntervalString != "" && !c.Config.StoreIntervalIsSet {
		storeInterval, err := util.CutSeconds(JSONConfig.StoreIntervalString)
		if err != nil {
			log.Fatal(err)
		}
		c.WithStoreInterval(storeInterval)
	}

	if JSONConfig.CryptoKey != defaultCryptoKey && !c.Config.CryptoKeyIsSet {
		c.WithCryptoKey(JSONConfig.CryptoKey)
	}

	if !JSONConfig.RestoreEnable && defaultRestoreEnable && !c.Config.RestoreEnvIsSet { //nolint:all
		c.WithRestoreEnable(JSONConfig.RestoreEnable)
	}
	return c
}

// FromEnv populates the ServerConfig from environment variables.
func (c *ServerConfigBuilder) FromEnv() *ServerConfigBuilder {
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
	_, storeIntSet := os.LookupEnv("STORE_INTERVAL")
	if storeIntSet {
		c.Config.StoreIntervalIsSet = true
	}
	_, pathSet := os.LookupEnv("FILE_STORAGE_PATH")
	if pathSet {
		c.Config.StoreFilePathIsSet = true
	}
	_, restoreSet := os.LookupEnv("RESTORE")
	if restoreSet {
		c.Config.RestoreEnvIsSet = true
	}
	_, DSNSet := os.LookupEnv("DATABASE_DSN")
	if DSNSet {
		c.Config.DBAddressIsSet = true
	}
	return c
}

// GetConfigs returns the fully constructed ServerConfig by combining
// configurations from environment variables and command-line flags.
// It validates the address to ensure it is properly set.
func GetConfigs() (ServerConfig, error) {
	var c ServerConfigBuilder
	c.Config.SetDefaults()
	c.FromEnv().FromFlags().FromFile()
	if !util.ValidateAddress(c.Config.Address) {
		return ServerConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}

func GetTestConfig() (ServerConfig, error) {
	var c ServerConfigBuilder
	c.WithRestoreEnable(false).WithAddress("localhost:8080").WithStoreFilePath("./tmp/metrics-test.json")
	if !util.ValidateAddress(c.Config.Address) {
		return ServerConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}
