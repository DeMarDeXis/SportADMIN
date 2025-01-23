package postgres

import (
	"AdminAppForDiplom/internal/models/nhl"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type LoaderStorage struct {
	db *sqlx.DB
}

func NewLoaderStorage(db *sqlx.DB) *LoaderStorage {
	return &LoaderStorage{
		db: db,
	}
}

// TODO: fix func
func (l *LoaderStorage) Loader(teams []nhl.TeamDB) error {
	q := `INSERT INTO teams (name, abbreviation) VALUES (:name, :abbreviation)`

	for _, team := range teams {
		_, err := l.db.NamedExec(q, map[string]interface{}{
			"name":         team.Name,
			"abbreviation": team.Abbreviation,
		})
		if err != nil {
			return fmt.Errorf("failed to insert team %s: %w", team.Name, err)
		}
	}

	return nil
}
