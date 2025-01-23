package service

import (
	"AdminAppForDiplom/internal/models/nhl"
	"AdminAppForDiplom/internal/storage"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type LoadService struct {
	storage storage.Load
	log     *slog.Logger
}

func NewLoadService(storage storage.Load, log *slog.Logger) *LoadService {
	return &LoadService{
		storage: storage,
		log:     log,
	}
}
func (l *LoadService) Loader() error {
	teams, err := readTeamsFromJSON("./jsondata/NHLAbbr.json")
	if err != nil {
		return fmt.Errorf("failed to read teams: %w", err)
	}

	return l.storage.Loader(teams)
}

func readTeamsFromJSON(filepath string) ([]nhl.TeamDB, error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var teams []nhl.TeamDB
	if err := json.Unmarshal(fileData, &teams); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return teams, nil
}
