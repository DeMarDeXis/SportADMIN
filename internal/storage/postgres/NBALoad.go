package postgres

import (
	"AdminAppForDiplom/internal/models/nba"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type NBALoader struct {
	db *sqlx.DB
}

func NewNBALoader(db *sqlx.DB) *NBALoader {
	return &NBALoader{db: db}
}

func (l NBALoader) AbbrNBALoader(abbrTeams []nba.TeamDB) error {
	tx, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, abbrTeam := range abbrTeams {
		var teamID int
		q1 := `INSERT INTO teams (name, abbreviation, img_url)
				VALUES ($1, $2, $3)
				RETURNING id`
		err := tx.QueryRow(q1, abbrTeam.Name, abbrTeam.Abbr, abbrTeam.Image).Scan(&teamID)
		if err != nil {
			return fmt.Errorf("failed to insert team %s: %w", abbrTeam.Name, err)
		}

		q2 := `INSERT INTO nba_teams (id_team, conference, division)
				VALUES ($1, $2, $3)`
		_, err = tx.Exec(q2, teamID, abbrTeam.Conference, abbrTeam.Divisions)
		if err != nil {
			return fmt.Errorf("failed to insert nba_teams %s: %w", abbrTeam.Name, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
