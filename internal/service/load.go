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
	storage storage.NHLLoadDB
	log     *slog.Logger
}

func NewLoadService(storage storage.NHLLoadDB, log *slog.Logger) *LoadService {
	return &LoadService{
		storage: storage,
		log:     log,
	}
}
func (l *LoadService) AbbrLoader() error {
	teams, err := readTeamsFromJSON("./jsondata/nhl/NHLAbbrCustom.json")
	if err != nil {
		return fmt.Errorf("failed to read teams: %w", err)
	}

	return l.storage.NHLLoader(teams)
}

func (l *LoadService) RosterLoader() error {
	fileData, err := os.ReadFile("./jsondata/nhl/NHLAllRoster.json")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var teamRosters []nhl.TeamRosterDB
	if err := json.Unmarshal(fileData, &teamRosters); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	for _, roster := range teamRosters {
		if err := l.storage.RosterLoaderToDB(roster); err != nil {
			return fmt.Errorf("failed to load roster for team %s: %w", roster.TeamName, err)
		}
	}

	return nil
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
