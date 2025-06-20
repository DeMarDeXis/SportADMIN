package nhl

import (
	"strings"
	"time"
)

type ScheduleJSON struct {
	Date    string `json:"Date" db:"date_game"`
	Time    string `json:"Time" db:"time_game"`
	Visitor string `json:"Visitor" db:"visitor_team_id"`
	G       []int  `json:"G"`
	Home    string `json:"Home" db:"home_team_id"`
	FIELD7  string `json:"FIELD7"`
	Att     int    `json:"Att." db:"attendance"`
	LOG     string `json:"LOG" db:"game_duration"`
	Notes   string `json:"Notes" db:"notes"`
}

type Schedule struct {
	Date         string `db:"date_game"`
	Time         string `db:"time_game"`
	VisitorTeam  string
	HomeTeam     string
	VisitorScore int    `db:"visitor_score"`
	HomeScore    int    `db:"home_score"`
	Attendance   int    `db:"attendance"`
	GameDuration string `db:"game_duration"`
	IsOvertime   bool   `db:"is_overtime"`
	Venue        string `db:"venue"`
	Notes        string `db:"notes"`
}

type ScheduleExport struct {
	ID           int       `db:"id"`
	Date         string    `db:"date_game"`
	Time         string    `db:"time_game"`
	VisitorTeam  string    `db:"visitor_team"`
	HomeTeam     string    `db:"home_team"`
	VisitorScore *int      `db:"visitor_score"`
	HomeScore    *int      `db:"home_score"`
	Attendance   *int      `db:"attendance"`
	GameDuration *string   `db:"game_duration"`
	IsOvertime   bool      `db:"is_overtime"`
	Venue        *string   `db:"venue"`
	Notes        *string   `db:"notes"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type ScheduleImport struct {
	//ID           int     `db:"id"`
	Date         string  `db:"date_game"`
	Time         string  `db:"time_game"`
	VisitorTeam  string  `db:"visitor_team"`
	HomeTeam     string  `db:"home_team"`
	VisitorScore *int    `db:"visitor_score"`
	HomeScore    *int    `db:"home_score"`
	Attendance   *int    `db:"attendance"`
	GameDuration *string `db:"game_duration"`
	IsOvertime   bool    `db:"is_overtime"`
	Venue        *string `db:"venue"`
	Notes        *string `db:"notes"`
}

func ParseScheduleFromJSON(jsonData ScheduleJSON) Schedule {
	var schedule Schedule
	schedule = Schedule{
		Date:         jsonData.Date,
		Time:         jsonData.Time,
		VisitorTeam:  jsonData.Visitor,
		HomeTeam:     jsonData.Home,
		Attendance:   jsonData.Att,
		GameDuration: jsonData.LOG,
		Notes:        jsonData.Notes,
	}

	if len(jsonData.G) >= 2 {
		schedule.VisitorScore = jsonData.G[0]
		schedule.HomeScore = jsonData.G[1]
	}

	schedule.IsOvertime = strings.Contains(strings.ToUpper(jsonData.FIELD7), "OT")

	return schedule
}
