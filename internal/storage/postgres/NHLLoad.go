package postgres

import (
	"AdminAppForDiplom/internal/models/nhl"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
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

func (l *LoaderStorage) ScheduleLoaderToDB(schedules []nhl.Schedule) error {
	const op = "internal/storage/postgres/NHLLoad.go"

	tx, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w \n in %s", err, op)
	}
	defer tx.Rollback()

	getTeamIDByNameStmt, err := tx.Prepare(`SELECT nhl_teams.id FROM nhl_teams
			JOIN teams ON teams.id = nhl_teams.id_team
			WHERE teams.name = $1`)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w \n in %s", err, op)
	}
	defer getTeamIDByNameStmt.Close()

	insertStmt, err := tx.Prepare(`INSERT INTO nhl_schedule (date_game,
                          time_game,
                          visitor_team_id,
                          home_team_id,
                          visitor_score,
                          home_score,
                          attendance,
                          game_duration,
                          is_overtime,
                          notes)
                          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w \n in %s", err, op)
	}
	defer insertStmt.Close()

	for _, schedule := range schedules {
		fmt.Printf("Parsing date: %s\n", schedule.Date)
		date, err := time.Parse("2006-01-02", schedule.Date)
		if err != nil {
			return fmt.Errorf("failed to parse date %s: %w", schedule.Date, err)
		}

		timeObj, err := time.Parse("3:04 PM", schedule.Time)
		if err != nil {
			return fmt.Errorf("failed to parse time %s: %w", schedule.Time, err)
		}

		var visitorTeamID int
		err = getTeamIDByNameStmt.QueryRow(schedule.VisitorTeam).Scan(&visitorTeamID)
		if err != nil {
			return fmt.Errorf("failed to get visitor team ID for %s: %w", schedule.VisitorTeam, err)
		}

		var homeTeamID int
		err = getTeamIDByNameStmt.QueryRow(schedule.HomeTeam).Scan(&homeTeamID)
		if err != nil {
			return fmt.Errorf("failed to get home team ID for %s: %w", schedule.HomeTeam, err)
		}

		_, err = insertStmt.Exec(date, timeObj,
			visitorTeamID, homeTeamID, schedule.VisitorScore, schedule.HomeScore,
			schedule.Attendance, schedule.GameDuration, schedule.IsOvertime, schedule.Notes)
		if err != nil {
			return fmt.Errorf("failed to insert schedule: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w \n in %s", err, op)
	}

	return nil
}

func (l *LoaderStorage) GetAllSchedule() ([]nhl.ScheduleExport, error) {
	const op = "internal/storage/postgres/NHLLoad.go"

	q := `SELECT
		ns.id,
		ns.date_game,
		ns.time_game,
		vt.name AS visitor_team,
		ht.name AS home_team,
		ns.visitor_score,
		ns.home_score,
		ns.attendance,
		ns.game_duration,
		ns.is_overtime,
		ns.venue,
		ns.notes
		FROM nhl_schedule ns
		JOIN
			nhl_teams nvt ON nvt.id = ns.visitor_team_id
		JOIN
			teams vt ON vt.id = nvt.id_team
		JOIN
			nhl_teams nht ON nht.id = ns.home_team_id
		JOIN
			teams ht ON ht.id = nht.id_team
			ORDER BY ns.date_game DESC`

	var schedules []nhl.ScheduleExport
	err := l.db.Select(&schedules, q)
	if err != nil {
		return nil, fmt.Errorf("failed to select schedules: %w \n in %s", err, op)
	}

	return schedules, nil
}

func (l *LoaderStorage) UpsertSchedule(schedules []nhl.ScheduleImport) error {
	const op = "internal/storage/postgres/NHLLoad.go"

	tx, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w \n in %s", err, op)
	}
	defer tx.Rollback()

	//Allocate a map to store team names and their corresponding IDs
	teamNames := make([]string, 0, len(schedules)*2)
	for _, s := range schedules {
		teamNames = append(teamNames, s.VisitorTeam, s.HomeTeam)
	}
	teamNames = uniqueStrings(teamNames) //filter unique team names

	teamIDs := make(map[string]int)
	if len(teamNames) > 0 {
		query, args, err := sqlx.In(`
            SELECT t.name, nt.id 
            FROM nhl_teams nt
            JOIN teams t ON t.id = nt.id_team
            WHERE t.name IN (?)`, teamNames)
		if err != nil {
			return fmt.Errorf("failed to build team query: %w \n in %s", err, op)
		}

		query = l.db.Rebind(query)
		rows, err := tx.Query(query, args...)
		if err != nil {
			return fmt.Errorf("failed to get team IDs: %w \n in %s", err, op)
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			var id int
			if err := rows.Scan(&name, &id); err != nil {
				return fmt.Errorf("failed to scan team ID: %w \n in %s", err, op)
			}
			teamIDs[name] = id
		}
	}

	//stmt, err := tx.Prepare(`
	//    INSERT INTO nhl_schedule (
	//        date_game, time_game, visitor_team_id, home_team_id,
	//        visitor_score, home_score, attendance, game_duration,
	//        is_overtime, venue, notes
	//    )
	//    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	//    ON CONFLICT (date_game, time_game, visitor_team_id, home_team_id)
	//    DO UPDATE SET
	//        visitor_score = EXCLUDED.visitor_score,
	//        home_score = EXCLUDED.home_score,
	//        attendance = EXCLUDED.attendance,
	//        game_duration = EXCLUDED.game_duration,
	//        is_overtime = EXCLUDED.is_overtime,
	//        venue = EXCLUDED.venue,
	//        notes = EXCLUDED.notes,
	//        updated_at = NOW()`)
	stmt, err := tx.Prepare(`
    INSERT INTO nhl_schedule (
        date_game, time_game, visitor_team_id, home_team_id,
        visitor_score, home_score, attendance, game_duration,
        is_overtime, venue, notes
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    ON CONFLICT (date_game, time_game, visitor_team_id, home_team_id) 
    DO UPDATE SET
        visitor_score = EXCLUDED.visitor_score,
        home_score = EXCLUDED.home_score,
        attendance = EXCLUDED.attendance,
        game_duration = EXCLUDED.game_duration,
        is_overtime = EXCLUDED.is_overtime,
        venue = EXCLUDED.venue,
        notes = EXCLUDED.notes,
        updated_at = CASE
            WHEN nhl_schedule.visitor_score IS DISTINCT FROM EXCLUDED.visitor_score OR
                 nhl_schedule.home_score IS DISTINCT FROM EXCLUDED.home_score OR
                 nhl_schedule.attendance IS DISTINCT FROM EXCLUDED.attendance OR
                 nhl_schedule.game_duration IS DISTINCT FROM EXCLUDED.game_duration OR
                 nhl_schedule.is_overtime IS DISTINCT FROM EXCLUDED.is_overtime OR
                 nhl_schedule.venue IS DISTINCT FROM EXCLUDED.venue OR
                 nhl_schedule.notes IS DISTINCT FROM EXCLUDED.notes
            THEN NOW()
            ELSE nhl_schedule.updated_at
        END`)
	if err != nil {
		return fmt.Errorf("failed to prepare upsert statement: %w \n in %s", err, op)
	}
	defer stmt.Close()

	for _, s := range schedules {
		visitorID, ok := teamIDs[s.VisitorTeam]
		if !ok {
			return fmt.Errorf("visitor team not found: %s \n in %s", s.VisitorTeam, op)
		}

		homeID, ok := teamIDs[s.HomeTeam]
		if !ok {
			return fmt.Errorf("home team not found: %s \n in %s", s.HomeTeam, op)
		}

		_, err := stmt.Exec(
			s.Date, s.Time, visitorID, homeID,
			s.VisitorScore, s.HomeScore, s.Attendance, s.GameDuration,
			s.IsOvertime, s.Venue, s.Notes)
		if err != nil {
			return fmt.Errorf("failed to upsert schedule: %w \n in %s", err, op)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w \n in %s", err, op)
	}

	return nil
}

// Вспомогательная функция для уникальных строк
func uniqueStrings(input []string) []string {
	unique := make(map[string]struct{})
	for _, v := range input {
		unique[v] = struct{}{}
	}
	result := make([]string, 0, len(unique))
	for k := range unique {
		result = append(result, k)
	}
	return result
}
