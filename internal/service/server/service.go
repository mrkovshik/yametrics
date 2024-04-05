package service

import (
	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	Storage storage.Storage
	Config  config.ServerConfig
	Logger  *zap.SugaredLogger
}

func NewServer(storage storage.Storage, cfg config.ServerConfig, logger *zap.SugaredLogger) *Server {
	return &Server{
		Storage: storage,
		Config:  cfg,
		Logger:  logger,
	}
}
