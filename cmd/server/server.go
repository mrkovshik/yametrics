package main

import (
	"github.com/mrkovshik/yametrics/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	service "github.com/mrkovshik/yametrics/internal/service/server"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("zap.NewDevelopment",
			zap.Error(err))
	}
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()

	mapStorage := storage.NewMapStorage()
	cfg, err := config.GetConfigs()
	if err != nil {
		sugar.Fatal("cfg.GetConfigs", err)
	}
	getMetricsService := service.NewServer(mapStorage, cfg, sugar)
	if cfg.RestoreEnable {
		if err := getMetricsService.RestoreMetrics(cfg.StoreFilePath); err != nil {
			sugar.Fatal("RestoreMetrics", err)
		}
	}
	if cfg.StoreEnable && !cfg.SyncStoreEnable {
		go getMetricsService.DumpMetrics()
	}
	run(getMetricsService, sugar, cfg)
	if err := getMetricsService.StoreMetrics(cfg.StoreFilePath); err != nil {
		sugar.Fatal("StoreMetrics", err)
	}

}

func run(s *service.Server, logger *zap.SugaredLogger, cfg config.ServerConfig) {
	r := chi.NewRouter()
	r.Use(s.WithLogging, s.GzipHandle)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", api.UpdateMetricFromJSONHandler(s))
		r.Post("/{type}/{name}/{value}", api.UpdateMetricFromURLHandler(s))
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", api.GetMetricFromJSONHandler(s))
		r.Get("/{type}/{name}", api.GetMetricFromURLHandler(s))
	})
	r.Get("/", api.GetMetricsHandler(s))
	logger.Infof("Starting server on %v\n", cfg.Address)
	logger.Fatal(http.ListenAndServe(cfg.Address, r))
}
