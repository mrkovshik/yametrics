package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	service "github.com/mrkovshik/yametrics/internal/service/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
	"log"
	"net/http"
)

func main() {
	cfg := config.ServerConfig{}
	mapStorage := server.NewMapStorage()
	if err := cfg.GetConfigs(); err != nil {
		log.Fatal(err)
	}
	getMetricsService := service.NewServer(mapStorage, log.Default(), cfg)
	run(getMetricsService)

}

func run(s *service.Server) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric(s))
	r.Get("/value/{type}/{name}", api.GetMetric(s))
	r.Get("/", api.GetMetrics(s))
	log.Printf("Starting server on %v\n", s.Config.Address)
	log.Fatal(http.ListenAndServe(s.Config.Address, r))
}
