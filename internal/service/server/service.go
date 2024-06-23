package service

import (
	"database/sql"

	"github.com/mrkovshik/yametrics/internal/service"
	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/server"
)

// Server represents the main application server containing configurations,
// logger, database connection, and storage service interface.
type Server struct {
	storage service.Storage     // storage service for handling metrics data
	config  config.ServerConfig // server configuration
	logger  *zap.SugaredLogger  // logger instance for logging
	db      *sql.DB             // database connection
}

// NewServer creates a new instance of Server with the provided dependencies.
// It initializes the server with the given storage service, configuration,
// logger, and database connection.
func NewServer(storage service.Storage, cfg config.ServerConfig, logger *zap.SugaredLogger, db *sql.DB) *Server {
	return &Server{
		db:      db,
		storage: storage,
		config:  cfg,
		logger:  logger,
	}
}
