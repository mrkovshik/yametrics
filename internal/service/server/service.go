package service

import (
	"database/sql"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	storage storage.Storage
	config  config.ServerConfig
	logger  *zap.SugaredLogger
	db      *sql.DB
}

func NewServer(storage storage.Storage, cfg config.ServerConfig, logger *zap.SugaredLogger, db *sql.DB) *Server {
	return &Server{
		db:      db,
		storage: storage,
		config:  cfg,
		logger:  logger,
	}
}
