package service

import (
	"log"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
)

type Server struct {
	Storage server.IStorage
	Logger  *log.Logger
	Config  config.ServerConfig
}

func NewServer(storage server.IStorage, logger *log.Logger, cfg config.ServerConfig) *Server {
	return &Server{
		Storage: storage,
		Logger:  logger,
		Config:  cfg,
	}
}
