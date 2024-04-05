package service

import (
	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	storage storage.Storage
	config  config.ServerConfig
	logger  *zap.SugaredLogger
}

func NewServer(storage storage.Storage, cfg config.ServerConfig, logger *zap.SugaredLogger) *Server {
	return &Server{
		storage: storage,
		config:  cfg,
		logger:  logger,
	}
}
