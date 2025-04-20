package postgres

import (
	"AdminAppForDiplom/internal/models/nfl"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type NFLLoader struct {
	db *sqlx.DB
}

func NewNFLLoader(db *sqlx.DB) *NFLLoader {
	return &NFLLoader{db: db}
}

func (n *NFLLoader) AbbrNFLLoader(abbrTeams []nfl.TeamDB) error {
	const pathError = "storage.NFLLoader.AbbrNFLLoader"

	tx, err := n.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", pathError, err)
	}
	defer tx.Rollback()

	for _, abbrTeam := range abbrTeams {
		var teamID int
		q1 := `INSERT INTO teams (name, abbreviation, img_url)
				VALUES ($1, $2, $3)
				RETURNING id`

		err := tx.QueryRow(q1, abbrTeam.Name, abbrTeam.Abbr, abbrTeam.ImgUrl).Scan(&teamID)
		if err != nil {
			return fmt.Errorf("failed to insert team %s: %w\n PATH ERROR: %s", abbrTeam.Name, err, pathError)
		}

		q2 := `INSERT INTO nfl_teams (id_team,com_used_abbr, conference, division) 
				VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(q2, teamID, abbrTeam.ComUsedAbbr, abbrTeam.Conference, abbrTeam.Divisions)
		if err != nil {
			return fmt.Errorf("failed to insert nfl_teams %s: %w\n PATH ERROR: %s", abbrTeam.Name, err, pathError)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
