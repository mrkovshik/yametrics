package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/mrkovshik/yametrics/internal/utl"
)

type ServerConfig struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	SyncStore     bool
	StoreFilePath string `env:"FILE_STORAGE_PATH"`
	StoreEnable   bool
	RestoreEnable bool `env:"RESTORE"`
	RestoreEnvSet bool
}

type ServerConfigBuilder struct {
	Config ServerConfig
}

func (c *ServerConfigBuilder) WithAddress(host string) *ServerConfigBuilder {
	c.Config.Address = host
	return c
}

func (c *ServerConfigBuilder) WithStoreInterval(interval int) *ServerConfigBuilder {
	c.Config.StoreInterval = interval
	return c
}
func (c *ServerConfigBuilder) WithStoreFilePath(path string) *ServerConfigBuilder {
	c.Config.StoreFilePath = path
	return c
}
func (c *ServerConfigBuilder) WithRestoreEnable(restore bool) *ServerConfigBuilder {
	c.Config.RestoreEnable = restore
	return c
}

func (c *ServerConfigBuilder) FromFlags() *ServerConfigBuilder {
	addr := flag.String("a", "localhost:8080", "server host and port")
	storeInterval := flag.Int("i", 300, "time interval between storing data to file")
	storeFilePath := flag.String("f", "/tmp/metrics-db.json", "path to storing data file")
	restoreEnable := flag.Bool("r", true, "is data restore from file enabled")

	flag.Parse()

	if c.Config.Address == "" {
		c.WithAddress(*addr)
	}
	if c.Config.StoreInterval == 0 && !c.Config.SyncStore {
		c.WithStoreInterval(*storeInterval)
	}
	if c.Config.StoreFilePath == "" && !c.Config.StoreEnable {
		c.WithStoreFilePath(*storeFilePath)
	}
	if !c.Config.RestoreEnvSet {
		c.WithRestoreEnable(*restoreEnable)
	}
	return c
}

func (c *ServerConfigBuilder) FromEnv() *ServerConfigBuilder {
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}
	_, storeIntSet := os.LookupEnv("STORE_INTERVAL")
	if storeIntSet {
		c.Config.SyncStore = true
	}
	_, pathSet := os.LookupEnv("FILE_STORAGE_PATH")
	if pathSet {
		c.Config.StoreEnable = true
	}
	_, restore := os.LookupEnv("RESTORE")
	if restore {
		c.Config.RestoreEnvSet = true
	}
	return c
}

func GetConfigs() (ServerConfig, error) {
	var c ServerConfigBuilder
	c.FromEnv().FromFlags()
	if !utl.ValidateAddress(c.Config.Address) {
		return ServerConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}
