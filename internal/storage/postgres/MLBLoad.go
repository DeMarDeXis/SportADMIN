package postgres

import (
	"AdminAppForDiplom/internal/models/mlb"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type MLBLoaderStorage struct {
	db *sqlx.DB
}

func NewMLBLoaderStorage(db *sqlx.DB) *MLBLoaderStorage {
	return &MLBLoaderStorage{
		db: db,
	}
}
func (l *MLBLoaderStorage) AbbrMLBLoader(teams []mlb.TeamDB) error {
	tx, err := l.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qForTeams := `INSERT INTO teams (name, abbreviation, img_url) VALUES ($1, $2, $3) RETURNING id`
	qForMLBTeams := `INSERT INTO mlb_teams (id_team, conference, division) VALUES ($1, $2, $3)`

	for _, team := range teams {
		var teamID int

		err := tx.QueryRow(qForTeams, team.Name, team.Abbr, team.Image).Scan(&teamID)
		if err != nil {
			return fmt.Errorf("error loading team %d: %v", teamID, err)
		}
		_, err = tx.Exec(qForMLBTeams, teamID, team.Conference, team.Division)
		if err != nil {
			return fmt.Errorf("error loading team %d: %v", teamID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
