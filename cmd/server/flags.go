package main

import (
	"flag"
	"github.com/mrkovshik/yametrics/internal/flags"
)

var addr = flags.NetAddress{
	Host: "localhost",
	Port: 8080,
}

func parseFlags() {
	flag.Var(&addr, "a", "address and port to run server")
	flag.Parse()
}
