package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	"go.uber.org/zap"
)

type restAPIServer struct {
	service api.Service
	config  *config.ServerConfig
	logger  *zap.SugaredLogger
}

func NewRestAPIServer(service api.Service, config *config.ServerConfig, logger *zap.SugaredLogger) api.Server {
	return &restAPIServer{
		service: service,
		config:  config,
		logger:  logger,
	}
}
func (s *restAPIServer) RunServer(ctx context.Context) {
	r := chi.NewRouter()
	r.Use(s.WithLogging, s.GzipHandle, s.Authenticate, s.SignResponse)
	r.Route("/update", func(r chi.Router) {
		r.Post("/", s.UpdateMetricFromJSON(ctx))
		r.Post("/{type}/{name}/{value}", s.UpdateMetricFromURL(ctx))
	})
	r.Post("/updates/", s.UpdateMetricsFromJSON(ctx))
	r.Route("/value", func(r chi.Router) {
		r.Post("/", s.GetMetricFromJSON(ctx))
		r.Get("/{type}/{name}", s.GetMetricFromURL(ctx))
	})

	r.Get("/ping", s.Ping(ctx))
	r.Get("/", s.GetMetrics(ctx))
	s.logger.Infof("Starting server on %v\n StoreInterval: %v\n"+
		"StoreIntervalSet: %v\nSyncStoreEnable: %v\nStoreFilePath: %v\nStoreFilePathSet: %v\n"+
		"StoreEnable: %v\nRestoreEnable: %v\nRestoreEnvSet: %v\nDBAddress: %v\nDBAddressIsSet: %v\nDBEnable: %v\n", s.config.Address, s.config.StoreInterval,
		s.config.StoreIntervalSet, s.config.SyncStoreEnable, s.config.StoreFilePath, s.config.StoreFilePathSet, s.config.StoreEnable,
		s.config.RestoreEnable, s.config.RestoreEnvSet, s.config.DBAddress, s.config.DBAddressIsSet, s.config.DBEnable)
	s.logger.Fatal(http.ListenAndServe(s.config.Address, r))
}
