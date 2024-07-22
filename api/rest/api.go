package rest

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

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
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		if err := s.server.ListenAndServe(); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-stop
		if err := s.server.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

// ConfigureRouter configures routes and middleware.
func (s *Server) ConfigureRouter() *Server {
	router := chi.NewRouter()
	router.Use(s.WithLogging, s.GzipHandle, s.SignResponse, s.DecryptRequest, s.Authenticate)
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

	s.logger.Infof(
		"Starting server on %v\n "+
			"StoreInterval: %v\n"+
			"StoreIntervalIsSet: %v\n"+
			"SyncStoreEnable: %v\n"+
			"StoreFilePath: %v\n"+
			"StoreFilePathSet: %v\n"+
			"StoreEnable: %v\n"+
			"RestoreEnable: %v\n"+
			"RestoreEnvSet: %v\n"+
			"DBAddress: %v\n"+
			"DBAddressIsSet: %v\n"+
			"DBEnable: %v\n"+
			"CryptoKey: %v\n"+
			"CryptoKeyIsSet: %v\n"+
			"Key: %v\n"+
			"KeyIsSet: %v\n"+
			"ConfigFilePath: %v\n"+
			"ConfigFilePathIsSet: %v\n",
		s.config.Address,
		s.config.StoreInterval,
		s.config.StoreIntervalIsSet,
		s.config.SyncStoreEnable,
		s.config.StoreFilePath,
		s.config.StoreFilePathIsSet,
		s.config.StoreEnable,
		s.config.RestoreEnable,
		s.config.RestoreEnvIsSet,
		s.config.DBAddress,
		s.config.DBAddressIsSet,
		s.config.DBEnable,
		s.config.CryptoKey,
		s.config.CryptoKeyIsSet,
		s.config.Key,
		s.config.KeyIsSet,
		s.config.ConfigFilePath,
		s.config.ConfigFilePathIsSet)
	s.server.Handler = router
	return s
}
