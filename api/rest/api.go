package rest

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
)

// Server represents the server configuration and dependencies.
type Server struct {
	service api.Service
	config  *config.ServerConfig
	logger  *zap.SugaredLogger
	router  *chi.Mux
}

// NewServer creates a new Server instance.
// Parameters:
// - service: an implementation of the api.Service interface.
// - config: server configuration settings.
// - logger: a sugared logger instance.
// Returns:
// - a pointer to the new Server instance.
func NewServer(service api.Service, config *config.ServerConfig, logger *zap.SugaredLogger) *Server {
	return &Server{
		service: service,
		config:  config,
		logger:  logger,
		router:  chi.NewRouter(),
	}
}

// RunServer starts the HTTP server with the configured router.
func (s *Server) RunServer(ctx context.Context) error {
	errCh := make(chan error)
	go func() {
		if err := http.ListenAndServe(s.config.Address, s.router); err != nil {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}

// ConfigureRouter configures routes and middleware.
func (s *Server) ConfigureRouter() {
	s.router.Use(s.WithLogging, s.GzipHandle, s.Authenticate, s.SignResponse)
	s.router.Route("/update", func(r chi.Router) {
		r.Post("/", s.UpdateMetricFromJSON)
		r.Post("/{type}/{name}/{value}", s.UpdateMetricFromURL)
	})
	s.router.Post("/updates/", s.UpdateMetricsFromJSON)
	s.router.Route("/value", func(r chi.Router) {
		r.Post("/", s.GetMetricFromJSON)
		r.Get("/{type}/{name}", s.GetMetricFromURL)
	})

	s.router.Get("/ping", s.Ping)
	s.router.Get("/", s.GetMetrics)
	s.logger.Infof("Starting server on %v\n StoreInterval: %v\n"+
		"StoreIntervalSet: %v\nSyncStoreEnable: %v\nStoreFilePath: %v\nStoreFilePathSet: %v\n"+
		"StoreEnable: %v\nRestoreEnable: %v\nRestoreEnvSet: %v\nDBAddress: %v\nDBAddressIsSet: %v\nDBEnable: %v\n", s.config.Address, s.config.StoreInterval,
		s.config.StoreIntervalSet, s.config.SyncStoreEnable, s.config.StoreFilePath, s.config.StoreFilePathSet, s.config.StoreEnable,
		s.config.RestoreEnable, s.config.RestoreEnvSet, s.config.DBAddress, s.config.DBAddressIsSet, s.config.DBEnable)
}
