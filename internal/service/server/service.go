package service

import (
	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
	"go.uber.org/zap"
)

type Server struct {
	Storage storage.IServerStorage
	Config  config.ServerConfig
	Logger  *zap.SugaredLogger
}

func NewServer(storage storage.IServerStorage, cfg config.ServerConfig, logger *zap.SugaredLogger) *Server {
	return &Server{
		Storage: storage,
		Config:  cfg,
		Logger:  logger,
	}
}
