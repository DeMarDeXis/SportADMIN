package service

import (
	"AdminAppForDiplom/internal/storage"
	"log/slog"
)

type NHLLoad interface {
	AbbrLoader() error
	RosterLoader() error
}

type Service struct {
	NHLLoad
}

func NewService(storage *storage.Storage, log *slog.Logger) *Service {
	return &Service{
		NHLLoad: NewLoadService(storage.NHLLoadDB, log),
	}
}
