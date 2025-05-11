package service

import (
	"AdminAppForDiplom/internal/storage"
	"log/slog"
)

type NHLLoad interface {
	AbbrLoader() error
	RosterLoader() error
	ScheduleLoader(filePath string) error
	ExportScheduleToExcel(filePath string) error
	ImportScheduleFromExcel(filePath string) error
	AddNewMatchDataFromExcel() error
}

type NBALoad interface {
	AbbrNBALoader() error
}

type NFLLoad interface {
	AbbrNFLLoader() error
}

type Service struct {
	NHLLoad
	NBALoad
	NFLLoad
}

func NewService(storage *storage.Storage, log *slog.Logger) *Service {
	return &Service{
		NHLLoad: NewNHLLoadService(storage.NHLLoadDB, log),
		NBALoad: NewNBALoadService(storage.NBALoadDB, log),
		NFLLoad: NewNFLLoadService(storage.NFLLoadDB, log),
	}
}
