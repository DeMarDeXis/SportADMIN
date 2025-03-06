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

func (l *LoaderStorage) NHLLoader(teams []nhl.TeamDB) error {
	tx, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, team := range teams {
		var teamID int
		q1 := `INSERT INTO teams (name, abbreviation, img_url) 
               VALUES ($1, $2, $3) 
               RETURNING id`
		err := tx.QueryRow(q1, team.Name, team.Abbreviation, team.ImgURL).Scan(&teamID)
		if err != nil {
			return fmt.Errorf("failed to insert team %s: %w", team.Name, err)
		}

		q2 := `INSERT INTO nhl_teams (id_team, conference, division) 
               VALUES ($1, $2, $3)`
		_, err = tx.Exec(q2, teamID, team.Conference, team.Division)
		if err != nil {
			return fmt.Errorf("failed to insert nhl team %s: %w", team.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (l *LoaderStorage) RosterLoaderToDB(teamRoster nhl.TeamRosterDB) error {
	tx, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Get team_id from nhl_teams by team name
	var teamID int
	q := `SELECT nt.id FROM nhl_teams nt 
          JOIN teams t ON t.id = nt.id_team 
          WHERE t.name = $1`
	err = tx.QueryRow(q, teamRoster.TeamName).Scan(&teamID)
	if err != nil {
		return fmt.Errorf("failed to get team id for %s: %w", teamRoster.TeamName, err)
	}

	// Insert players
	for _, player := range teamRoster.Players {
		q := `INSERT INTO nhl_roster 
              (id_team, name, surname, number, position, hand, age, 
               acquired_at, birth_place, role, injured)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

		_, err = tx.Exec(q,
			teamID,
			player.Name,
			player.Surname,
			player.Number,
			player.Position,
			player.Hand,
			player.Age,
			player.Acquired,
			player.Birthplace,
			player.Role,
			player.Injured)

		if err != nil {
			return fmt.Errorf("failed to insert player %s %s: %w",
				player.Name, player.Surname, err)
		}
	}

	return tx.Commit()
}
