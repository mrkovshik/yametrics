package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrkovshik/yametrics/api"
	"log"
	"net/http"
	"time"

	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/storage"
)

const (
	readTimeout  = 5 * time.Second  // Adjust as needed
	writeTimeout = 10 * time.Second // Adjust as needed
	idleTimeout  = 15 * time.Second // Adjust as needed
)

func main() {
	mapStorage := storage.NewMapStorage()
	getMetricsService := service.NewServiceWithMapStorage(mapStorage, log.Default())
	run(getMetricsService)

}

func run(s *service.Service) {

	parseFlags()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric(s))
	r.Get("/value/{type}/{name}", api.GetMetric(s))
	r.Get("/", api.GetMetrics(s))
	fmt.Println("Running server on", addr.String())
	server := &http.Server{
		Addr:         addr.String(),
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
	log.Fatal(server.ListenAndServe())
}
