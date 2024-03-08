package main

import (
	"flag"
	"github.com/mrkovshik/yametrics/internal/flags"
	"time"
)

var addr = flags.NetAddress{
	Host: "localhost",
	Port: 8080,
}
var pollInterval time.Duration
var reportInterval time.Duration

func parseFlags() {

	flag.Var(&addr, "a", "address and port to run server")
	flag.DurationVar(&pollInterval, "r", 2*time.Second, "metrics polling interval")
	flag.DurationVar(&reportInterval, "p", 10*time.Second, "metrics sending to server interval")
	flag.Parse()
}
