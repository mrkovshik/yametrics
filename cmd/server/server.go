package main

import (
	"errors"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/storage"
	"github.com/mrkovshik/yametrics/internal/utl"
	"log"
	"net/http"
)

var hostPort *string

func main() {
	mapStorage := storage.NewMapStorage()
	getMetricsService := service.NewServiceWithMapStorage(mapStorage, log.Default())
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}
	run(getMetricsService)

}

func run(s *service.Service) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric(s))
	r.Get("/value/{type}/{name}", api.GetMetric(s))
	r.Get("/", api.GetMetrics(s))
	log.Println("Starting server on", *hostPort)
	log.Fatal(http.ListenAndServe(*hostPort, r))
}

func parseFlags() error {
	var cfg config.AgentConfig
	hostPort = flag.String("a", "localhost:8080", "server host and port")
	flag.Parse()
	if !utl.ValidateAddress(*hostPort) {
		return errors.New("need address in a form host:port")
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Address != "" {
		if !utl.ValidateAddress(cfg.Address) {
			log.Fatal(errors.New("invalid address env"))
		}
		hostPort = &cfg.Address
	}

	return nil
}
