package config

import (
	"errors"
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/mrkovshik/yametrics/internal/utl"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

type ServerConfigBuilder struct {
	Config ServerConfig
}

func (c *ServerConfigBuilder) WithAddress(host string) *ServerConfigBuilder {
	c.Config.Address = host
	return c
}

func (c *ServerConfigBuilder) FromFlags() *ServerConfigBuilder {
	addr := flag.String("a", "localhost:8080", "server host and port")
	flag.Parse()

	if c.Config.Address == "" {
		c.WithAddress(*addr)
	}
	return c
}

func (c *ServerConfigBuilder) FromEnv() *ServerConfigBuilder {
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
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
