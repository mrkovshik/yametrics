package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"

	"github.com/mrkovshik/yametrics/internal/util"
)

type ServerConfig struct {
	Address          string `env:"ADDRESS"`
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
}

type ServerConfigBuilder struct {
	Config ServerConfig
}

func (c *ServerConfigBuilder) WithAddress(host string) *ServerConfigBuilder {
	c.Config.Address = host
	return c
}

func (c *ServerConfigBuilder) WithDSN(dsn string) *ServerConfigBuilder {
	c.Config.DBAddress = dsn
	return c
}
func (c *ServerConfigBuilder) WithDBEnable() *ServerConfigBuilder {
	c.Config.DBEnable = true
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
func (c *ServerConfigBuilder) WithStoreEnable(store bool) *ServerConfigBuilder {
	c.Config.StoreEnable = store
	return c
}
func (c *ServerConfigBuilder) WithSyncStoreEnable(sync bool) *ServerConfigBuilder {
	c.Config.SyncStoreEnable = sync
	return c
}

func (c *ServerConfigBuilder) FromFlags() *ServerConfigBuilder {
	addr := flag.String("a", "localhost:8080", "server host and port")
	storeInterval := flag.Int("i", 300, "time interval between storing data to file")
	storeFilePath := flag.String("f", "./tmp/metrics-db.json", "path to storing data file")
	restoreEnable := flag.Bool("r", true, "is data restore from file enabled")
	dbAddress := flag.String("d", "", "db address") //host=localhost port=5432 user=yandex password=yandex dbname=yandex sslmode=disable
	flag.Parse()

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
	if c.Config.StoreInterval == 0 && c.Config.StoreEnable { //Если функция записи в файл отключена, то и эта опция будет отключена
		c.WithSyncStoreEnable(true)
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
	return c
}

func GetConfigs() (ServerConfig, error) {
	var c ServerConfigBuilder
	c.FromEnv().FromFlags()
	if !util.ValidateAddress(c.Config.Address) {
		return ServerConfig{}, errors.New("need address in a form host:port")
	}
	return c.Config, nil
}
