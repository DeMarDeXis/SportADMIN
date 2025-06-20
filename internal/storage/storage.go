package storage

import (
	"AdminAppForDiplom/internal/models/mlb"
	"AdminAppForDiplom/internal/models/nba"
	"AdminAppForDiplom/internal/models/nfl"
	"AdminAppForDiplom/internal/models/nhl"
	"AdminAppForDiplom/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type NHLLoadDB interface {
	NHLLoader(teams []nhl.TeamDB) error
	RosterLoaderToDB(teamRoster nhl.TeamRosterDB) error
	ScheduleLoaderToDB(schedules []nhl.Schedule) error
	GetAllSchedule() ([]nhl.ScheduleExport, error)
	UpsertSchedule(schedules []nhl.ScheduleImport) error
	AddNewSchedule(schedules []nhl.ScheduleImport) error
}

type NBALoadDB interface {
	AbbrNBALoader(abbrTeams []nba.TeamDB) error
}

type NFLLoadDB interface {
	AbbrNFLLoader(abbrTeams []nfl.TeamDB) error
}

type MLBLoadDB interface {
	AbbrMLBLoader(teams []mlb.TeamDB) error
}

type Storage struct {
	NHLLoadDB
	NBALoadDB
	NFLLoadDB
	MLBLoadDB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		NHLLoadDB: postgres.NewLoaderStorage(db),
		NBALoadDB: postgres.NewNBALoader(db),
		NFLLoadDB: postgres.NewNFLLoader(db),
		MLBLoadDB: postgres.NewMLBLoaderStorage(db),
	}
}
