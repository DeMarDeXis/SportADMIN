package storage

import (
	"AdminAppForDiplom/internal/models/nhl"
	"AdminAppForDiplom/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type NHLLoadDB interface {
	NHLLoader(teams []nhl.TeamDB) error
	RosterLoaderToDB(teamRoster nhl.TeamRosterDB) error
}

type Storage struct {
	NHLLoadDB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		NHLLoadDB: postgres.NewLoaderStorage(db),
	}
}
