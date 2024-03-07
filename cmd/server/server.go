package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrkovshik/yametrics/api"
	"log"
	"net/http"

	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/storage"
)

func main() {
	mapStorage := storage.NewMapStorage()

	getMetricsService := service.NewServiceWithMapStorage(mapStorage, log.Default())
	run(getMetricsService)

}

func run(s *service.Service) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric(s))
	r.Get("/value/{type}/{name}", api.GetMetric(s))
	r.Get("/", api.GetMetrics(s))
	log.Fatal(http.ListenAndServe(":8080", r))
}
