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

	var teamID int
	q := `SELECT nt.id FROM nhl_teams nt 
          JOIN teams t ON t.id = nt.id_team 
          WHERE t.name = $1`
	err = tx.QueryRow(q, teamRoster.TeamName).Scan(&teamID)
	if err != nil {
		return fmt.Errorf("failed to get team id for %s: %w", teamRoster.TeamName, err)
	}

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

	teamNames := make([]string, 0, len(schedules)*2)
	for _, s := range schedules {
		teamNames = append(teamNames, s.VisitorTeam, s.HomeTeam)
	}
	teamNames = uniqueStrings(teamNames)

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

	insertStmt, err := tx.Prepare(`
		INSERT INTO nhl_schedule (
			date_game, time_game, visitor_team_id, home_team_id,
			visitor_score, home_score, attendance, game_duration,
			is_overtime, venue, notes
		)
		VALUES ($1, $2, $3, $4, $5::integer, $6::integer, $7::integer, $8, $9, $10, $11)
		ON CONFLICT (date_game, time_game, visitor_team_id, home_team_id) 
		DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w \n in %s", err, op)
	}
	defer insertStmt.Close()

	updateStmt, err := tx.Prepare(`
		UPDATE nhl_schedule
		SET 
			visitor_score = $5::integer,
			home_score = $6::integer,
			attendance = $7::integer,
			game_duration = $8,
			is_overtime = $9,
			venue = $10,
			notes = $11,
			updated_at = NOW()
		WHERE 
			date_game = $1 AND 
			time_game = $2 AND 
			visitor_team_id = $3 AND 
			home_team_id = $4 AND
			(
				visitor_score IS DISTINCT FROM $5::integer OR
				home_score IS DISTINCT FROM $6::integer OR
				attendance IS DISTINCT FROM $7::integer OR
				game_duration IS DISTINCT FROM $8 OR
				is_overtime IS DISTINCT FROM $9 OR
				venue IS DISTINCT FROM $10 OR
				notes IS DISTINCT FROM $11
			)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w \n in %s", err, op)
	}
	defer updateStmt.Close()

	for _, s := range schedules {
		visitorID, ok := teamIDs[s.VisitorTeam]
		if !ok {
			return fmt.Errorf("visitor team not found: %s \n in %s", s.VisitorTeam, op)
		}

		homeID, ok := teamIDs[s.HomeTeam]
		if !ok {
			return fmt.Errorf("home team not found: %s \n in %s", s.HomeTeam, op)
		}

		_, err := insertStmt.Exec(
			s.Date, s.Time, visitorID, homeID,
			s.VisitorScore, s.HomeScore, s.Attendance, s.GameDuration,
			s.IsOvertime, s.Venue, s.Notes)
		if err != nil {
			return fmt.Errorf("failed to insert schedule: %w \n in %s", err, op)
		}

		result, err := updateStmt.Exec(
			s.Date, s.Time, visitorID, homeID,
			s.VisitorScore, s.HomeScore, s.Attendance, s.GameDuration,
			s.IsOvertime, s.Venue, s.Notes)
		if err != nil {
			return fmt.Errorf("failed to update schedule: %w \n in %s", err, op)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			fmt.Printf("Updated %d rows for game %s %s: %s vs %s\n",
				rowsAffected, s.Date, s.Time, s.VisitorTeam, s.HomeTeam)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w \n in %s", err, op)
	}

	return nil
}

func (l *LoaderStorage) AddNewSchedule(schedules []nhl.ScheduleImport) error {
	const op = "storage.postgres.NHLLoad.AddNewSchedule"

	tx, err := l.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w \n in %s", err, op)
	}
	defer tx.Rollback()

	teamNames := make([]string, 0, len(schedules)*2)
	for _, s := range schedules {
		teamNames = append(teamNames, s.VisitorTeam, s.HomeTeam)
	}
	teamNames = uniqueStrings(teamNames)
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

	qNewSchedule, err := tx.Prepare(`INSERT INTO nhl_schedule (date_game, time_game, visitor_team_id, home_team_id, visitor_score, home_score, attendance, game_duration,is_overtime, venue, notes) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
					ON CONFLICT (date_game, time_game, visitor_team_id, home_team_id) DO NOTHING`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w \n in %s", err, op)
	}
	defer qNewSchedule.Close()

	for _, s := range schedules {
		visitorID, ok := teamIDs[s.VisitorTeam]
		if !ok {
			return fmt.Errorf("visitor team not found: %s \n in %s", s.VisitorTeam, op)
		}

		homeID, ok := teamIDs[s.HomeTeam]
		if !ok {
			return fmt.Errorf("home team not found: %s \n in %s", s.HomeTeam, op)
		}

		_, err := qNewSchedule.Exec(s.Date, s.Time, visitorID, homeID,
			s.VisitorScore, s.HomeScore, s.Attendance, s.GameDuration, s.IsOvertime, s.Venue, s.Notes)
		if err != nil {
			return fmt.Errorf("failed to insert schedule: %w \n in %s", err, op)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w \n in %s", err, op)
	}

	return nil
}

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
