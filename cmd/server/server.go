package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	service "github.com/mrkovshik/yametrics/internal/service/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	logger := httplog.NewLogger("httplog-example", httplog.Options{
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
	})

	mapStorage := storage.NewMapStorage()
	cfg, err := config.GetConfigs()
	if err != nil {
		logger.Error("cfg.GetConfigs", err)
	}
	getMetricsService := service.NewServer(mapStorage, cfg, logger)
	run(getMetricsService)

}

func run(s *service.Server) {
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(s.Logger))
	r.Post("/update/{type}/{name}/{value}", api.UpdateMetricHandler(s))
	r.Get("/value/{type}/{name}", api.GetMetricHandler(s))
	r.Get("/", api.GetMetricsHandler(s))
	log.Printf("Starting server on %v\n", s.Config.Address)
	log.Fatal(http.ListenAndServe(s.Config.Address, r))
}
