package rest

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mrkovshik/yametrics/api"
	config "github.com/mrkovshik/yametrics/internal/config/server"
)

// Server represents the server configuration and dependencies.
type Server struct {
	server  *http.Server
	service api.Service
	config  *config.ServerConfig
	logger  *zap.SugaredLogger
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
		server: &http.Server{
			Addr: config.Address,
		},
		service: service,
		config:  config,
		logger:  logger,
	}
}

// RunServer starts the HTTP server with the configured router.
func (s *Server) RunServer(stop chan os.Signal) error {
	errCh := make(chan error)
	go func(server *http.Server) {
		if err := server.ListenAndServe(); err != nil {
			errCh <- err
		}
	}(s.server)
	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-stop:
		err := s.server.Shutdown(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

// ConfigureRouter configures routes and middleware.
func (s *Server) ConfigureRouter() *Server {
	router := chi.NewRouter()
	router.Use(s.WithLogging, s.GzipHandle, s.Authenticate, s.SignResponse)
	router.Route("/update", func(r chi.Router) {
		r.Post("/", s.HandleUpdateMetricFromJSON)
		r.Post("/{type}/{name}/{value}", s.HandleUpdateMetricFromURL)
	})
	router.Post("/updates/", s.HandleUpdateMetricsFromJSON)
	router.Route("/value", func(r chi.Router) {
		r.Post("/", s.HandleGetMetricFromJSON)
		r.Get("/{type}/{name}", s.HandleGetMetricFromURL)
	})

	router.Get("/ping", s.HandlePing)
	router.Get("/", s.HandleGetMetrics)
	s.logger.Infof("Starting server on %v\n StoreInterval: %v\n"+
		"StoreIntervalSet: %v\nSyncStoreEnable: %v\nStoreFilePath: %v\nStoreFilePathSet: %v\n"+
		"StoreEnable: %v\nRestoreEnable: %v\nRestoreEnvSet: %v\nDBAddress: %v\nDBAddressIsSet: %v\nDBEnable: %v\n", s.config.Address, s.config.StoreInterval,
		s.config.StoreIntervalSet, s.config.SyncStoreEnable, s.config.StoreFilePath, s.config.StoreFilePathSet, s.config.StoreEnable,
		s.config.RestoreEnable, s.config.RestoreEnvSet, s.config.DBAddress, s.config.DBAddressIsSet, s.config.DBEnable)
	s.server.Handler = router
	return s
}
