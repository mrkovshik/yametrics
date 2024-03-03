package service

import (
	"log"

	"github.com/mrkovshik/yametrics/internal/storage"
)

type Service struct {
	Storage storage.IStorage
	Logger  *log.Logger
}

func NewServiceWithMapStorage(storage storage.IStorage, logger *log.Logger) *Service {
	return &Service{
		Storage: storage,
		Logger:  logger,
	}
}
