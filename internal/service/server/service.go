package service

import (
	"log"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
)

type Server struct {
	Storage storage.IServerStorage
	Config  config.ServerConfig
	Logger  *log.Logger
}

func NewServer(storage storage.IServerStorage, cfg config.ServerConfig, logger *log.Logger) *Server {
	return &Server{
		Storage: storage,
		Config:  cfg,
		Logger:  logger,
	}
}
