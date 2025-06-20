package service

import (
	"AdminAppForDiplom/internal/domain/mlbFiles"
	"AdminAppForDiplom/internal/models/mlb"
	"AdminAppForDiplom/internal/storage"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

type MLBLoadService struct {
	storage storage.MLBLoadDB
	log     *slog.Logger
}

func NewMLBLoadService(storage storage.MLBLoadDB, log *slog.Logger) *MLBLoadService {
	return &MLBLoadService{
		storage: storage,
		log:     log,
	}
}

func (s *MLBLoadService) AbbrMLBLoader() error {
	fileData, err := s.readMLBAbbrFromJSON(mlbFiles.MLBAbbrJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return s.storage.AbbrMLBLoader(fileData)
}

func (s *MLBLoadService) readMLBAbbrFromJSON(filepath string) ([]mlb.TeamDB, error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var teams []mlb.TeamDB
	if err := json.Unmarshal(fileData, &teams); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return teams, nil
}
