package service

import (
	"AdminAppForDiplom/internal/models/nfl"
	"AdminAppForDiplom/internal/storage"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type NFLLoadService struct {
	storage storage.NFLLoadDB
	log     *slog.Logger
}

func NewNFLLoadService(storage storage.NFLLoadDB, log *slog.Logger) *NFLLoadService {
	return &NFLLoadService{storage: storage, log: log}
}

func (l *NFLLoadService) AbbrNFLLoader() error {
	abbr, err := readNFLAbbrFromJSON("./jsondata/nfl/NFLAbbr.json")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return l.storage.AbbrNFLLoader(abbr)
}

func readNFLAbbrFromJSON(filepath string) ([]nfl.TeamDB, error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var abbrTeams []nfl.TeamDB
	if err := json.Unmarshal(fileData, &abbrTeams); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return abbrTeams, nil
}
