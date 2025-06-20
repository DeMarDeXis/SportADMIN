package service

import (
	"AdminAppForDiplom/internal/models/nba"
	"AdminAppForDiplom/internal/storage"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type NBALoadService struct {
	storage storage.NBALoadDB
	log     *slog.Logger
}

func NewNBALoadService(storage storage.NBALoadDB, log *slog.Logger) *NBALoadService {
	return &NBALoadService{
		storage: storage,
		log:     log,
	}

}

func (l *NBALoadService) AbbrNBALoader() error {
	abbr, err := readNBAAbbrFromJSON("./jsondata/nba/NBAAbbrCUSTOM.json")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return l.storage.AbbrNBALoader(abbr)
}

func readNBAAbbrFromJSON(filepath string) ([]nba.TeamDB, error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var abbrTeams []nba.TeamDB
	if err := json.Unmarshal(fileData, &abbrTeams); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return abbrTeams, nil
}
