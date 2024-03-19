package service

import (
	"github.com/go-chi/httplog/v2"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
)

type Server struct {
	Storage storage.IServerStorage
	Config  config.ServerConfig
	Logger  *httplog.Logger
}

func NewServer(storage storage.IServerStorage, cfg config.ServerConfig, logger *httplog.Logger) *Server {
	return &Server{
		Storage: storage,
		Config:  cfg,
		Logger:  logger,
	}
}
