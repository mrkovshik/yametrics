package agent

import (
	"errors"
	"flag"
	"github.com/mrkovshik/yametrics/internal/utl"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func (c *ServerConfig) GetConfigs() error {
	addr := flag.String("a", "localhost:8080", "server host and port")
	flag.Parse()
	c.Address = *addr
	if !utl.ValidateAddress(c.Address) {
		return errors.New("need address in a form host:port")
	}
	return nil
}
