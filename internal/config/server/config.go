package agent

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/mrkovshik/yametrics/internal/utl"
	"log"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func (c *ServerConfig) GetConfigs() error {
	addr := flag.String("a", "localhost:8080", "server host and port")
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
	return nil
}
