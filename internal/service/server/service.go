package server

import (
	"database/sql"

	"github.com/mrkovshik/yametrics/internal/service"
	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/server"
)

type Server struct {
	storage service.Storage
	config  config.ServerConfig
	Logger  *zap.SugaredLogger
	db      *sql.DB
}

func NewServer(storage service.Storage, cfg config.ServerConfig, logger *zap.SugaredLogger, db *sql.DB) *Server {
	return &Server{
		db:      db,
		storage: storage,
		config:  cfg,
		Logger:  logger,
	}
}
