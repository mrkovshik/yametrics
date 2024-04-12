package service

import (
	"database/sql"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/service"
	"go.uber.org/zap"
)

type Server struct {
	storage service.Storage
	config  config.ServerConfig
	logger  *zap.SugaredLogger
	db      *sql.DB
}

func NewServer(storage service.Storage, cfg config.ServerConfig, logger *zap.SugaredLogger, db *sql.DB) *Server {
	return &Server{
		db:      db,
		storage: storage,
		config:  cfg,
		logger:  logger,
	}
}
