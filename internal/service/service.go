package service

import (
	"AdminAppForDiplom/internal/storage"
	"log/slog"
)

type Load interface {
	Loader() error
}

type Service struct {
	Load
}

func NewService(storage *storage.Storage, log *slog.Logger) *Service {
	return &Service{
		Load: NewLoadService(storage.Load, log),
	}
}
