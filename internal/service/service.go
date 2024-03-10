package service

import (
	config "github.com/mrkovshik/yametrics/internal/config/server"
	"log"

	"github.com/mrkovshik/yametrics/internal/storage"
)

type Service struct {
	Storage storage.IStorage
	Logger  *log.Logger
	Config  config.ServerConfig
}

func NewServiceWithMapStorage(storage storage.IStorage, logger *log.Logger, cfg config.ServerConfig) *Service {
	return &Service{
		Storage: storage,
		Logger:  logger,
		Config:  cfg,
	}
}
