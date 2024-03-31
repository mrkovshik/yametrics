package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	service "github.com/mrkovshik/yametrics/internal/service/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
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
	if getMetricsService.Config.StoreEnable && !getMetricsService.Config.SyncStore {
		go getMetricsService.DumpMetrics()
	}
	run(getMetricsService)

}

func run(s *service.Server) {
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
	s.Logger.Infof("Starting server on %v\n", s.Config.Address)
	s.Logger.Fatal(http.ListenAndServe(s.Config.Address, r))
}
