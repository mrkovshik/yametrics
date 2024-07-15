// Package config provides configuration handling for the server, allowing
// configurations to be set via environment variables or command-line flags.
package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"

	"github.com/mrkovshik/yametrics/internal/util"
)

// ServerConfig holds the configuration settings for the server.
type ServerConfig struct {
	Address          string `env:"ADDRESS"`
	Key              string `env:"KEY"`
	StoreInterval    int    `env:"STORE_INTERVAL" envDefault:"300"`
	StoreIntervalSet bool
	SyncStoreEnable  bool   `envDefault:"false"`
	StoreFilePath    string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/metrics-db.json"`
	StoreFilePathSet bool
	StoreEnable      bool `envDefault:"true"`
	RestoreEnable    bool `env:"RESTORE" envDefault:"true"`
	RestoreEnvSet    bool
	DBAddress        string `env:"DATABASE_DSN"`
	DBAddressIsSet   bool
	DBEnable         bool
	CryptoKey        string `env:"CRYPTO_KEY"`
}

// ServerConfigBuilder is a builder for constructing a ServerConfig instance.
type ServerConfigBuilder struct {
	Config ServerConfig
}

// WithKey sets the key in the ServerConfig.
func (c *ServerConfigBuilder) WithKey(key string) *ServerConfigBuilder {
	c.Config.Key = key
	return c
}

// WithAddress sets the address in the ServerConfig.
func (c *ServerConfigBuilder) WithAddress(address string) *ServerConfigBuilder {
	c.Config.Address = address
	return c
}

// WithDSN sets the database DSN in the ServerConfig.
func (c *ServerConfigBuilder) WithDSN(dsn string) *ServerConfigBuilder {
	c.Config.DBAddress = dsn
	return c
}

// WithDBEnable enables the database in the ServerConfig.
func (c *ServerConfigBuilder) WithDBEnable() *ServerConfigBuilder {
	c.Config.DBEnable = true
	return c
}

// WithStoreInterval sets the store interval in the ServerConfig.
func (c *ServerConfigBuilder) WithStoreInterval(interval int) *ServerConfigBuilder {
	c.Config.StoreInterval = interval
	return c
}

// WithStoreFilePath sets the store file path in the ServerConfig.
func (c *ServerConfigBuilder) WithStoreFilePath(path string) *ServerConfigBuilder {
	c.Config.StoreFilePath = path
	return c
}

// WithRestoreEnable sets the restore enable flag in the ServerConfig.
func (c *ServerConfigBuilder) WithRestoreEnable(restore bool) *ServerConfigBuilder {
	c.Config.RestoreEnable = restore
	return c
}

// WithStoreEnable sets the store enable flag in the ServerConfig.
func (c *ServerConfigBuilder) WithStoreEnable(store bool) *ServerConfigBuilder {
	c.Config.StoreEnable = store
	return c
}

// WithSyncStoreEnable sets the sync store enable flag in the ServerConfig.
func (c *ServerConfigBuilder) WithSyncStoreEnable(sync bool) *ServerConfigBuilder {
	c.Config.SyncStoreEnable = sync
	return c
}

// WithCryptoKey sets the crypto key flag in the ServerConfig.
func (c *ServerConfigBuilder) WithCryptoKey(path string) *ServerConfigBuilder {
	c.Config.CryptoKey = path
	return c
}

// FromFlags populates the ServerConfig from command-line flags.
func (c *ServerConfigBuilder) FromFlags() *ServerConfigBuilder {
	addr := flag.String("a", "localhost:8080", "server host and port")
	storeInterval := flag.Int("i", 300, "time interval between storing data to file")
	storeFilePath := flag.String("f", "./tmp/metrics-db.json", "path to storing data file")
	restoreEnable := flag.Bool("r", true, "is data restore from file enabled")
	dbAddress := flag.String("d", "", "db address") //host=localhost port=5432 user=yandex password=yandex dbname=yandex sslmode=disable
	key := flag.String("k", "", "secret auth key")
	cryptoKey := flag.String("-crypto-key", "./tmp/private", "path to the file with private key")
	flag.Parse()

	if c.Config.Key == "" {
		c.WithKey(*key)
	}

	if c.Config.CryptoKey == "" {
		c.WithCryptoKey(*cryptoKey)
	}

	if c.Config.Address == "" {
		c.WithAddress(*addr)
	}
	if !c.Config.DBAddressIsSet {
		c.WithDSN(*dbAddress)
	}
	if c.Config.DBAddress != "" {
		c.WithDBEnable()
	}
	if !c.Config.StoreFilePathSet {
		c.WithStoreFilePath(*storeFilePath)
	}
	if c.Config.StoreFilePath == "" {
		c.WithStoreEnable(false)
	}
	if !c.Config.StoreIntervalSet {
		c.WithStoreInterval(*storeInterval)
	}
	if c.Config.StoreInterval == 0 && c.Config.StoreEnable { // If file storage is disabled, this option will be disabled as well
		c.WithSyncStoreEnable(true)
	}
	if !c.Config.RestoreEnvSet {
		c.WithRestoreEnable(*restoreEnable)
	}
	return c
}

// FromEnv populates the ServerConfig from environment variables.
func (c *ServerConfigBuilder) FromEnv() *ServerConfigBuilder {
	if err := env.Parse(&c.Config); err != nil {
		log.Fatal(err)
	}
	_, storeIntSet := os.LookupEnv("STORE_INTERVAL")
	if storeIntSet {
		c.Config.StoreIntervalSet = true
	}
	_, pathSet := os.LookupEnv("FILE_STORAGE_PATH")
	if pathSet {
		c.Config.StoreFilePathSet = true
	}
	_, restoreSet := os.LookupEnv("RESTORE")
	if restoreSet {
		c.Config.RestoreEnvSet = true
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
	c.FromEnv().FromFlags()
	if !util.ValidateAddress(c.Config.Address) {
		return ServerConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}

func GetTestConfig() (ServerConfig, error) {
	var c ServerConfigBuilder
	c.WithRestoreEnable(false).WithAddress("localhost:8080").WithStoreEnable(true).WithStoreFilePath("./tmp/metrics-test.json")
	if !util.ValidateAddress(c.Config.Address) {
		return ServerConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}
