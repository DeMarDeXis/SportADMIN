package storage

import (
	"AdminAppForDiplom/internal/models/nhl"
	"AdminAppForDiplom/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type Load interface {
	Loader(teams []nhl.TeamDB) error
}

type Storage struct {
	Load
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Load: postgres.NewLoaderStorage(db),
	}
}
