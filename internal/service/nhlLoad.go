package service

import (
	"AdminAppForDiplom/internal/models/nhl"
	"AdminAppForDiplom/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type NHLLoadService struct {
	storage storage.NHLLoadDB
	log     *slog.Logger
}

func NewNHLLoadService(storage storage.NHLLoadDB, log *slog.Logger) *NHLLoadService {
	return &NHLLoadService{
		storage: storage,
		log:     log,
	}
}
func (l *NHLLoadService) AbbrLoader() error {
	teams, err := readNHLAbbrFromJSON("./jsondata/nhl/NHLAbbrCustom.json")
	if err != nil {
		return fmt.Errorf("failed to read abbr: %w", err)
	}

	return l.storage.NHLLoader(teams)
}

func (l *NHLLoadService) RosterLoader() error {
	fileData, err := os.ReadFile("./jsondata/nhl/NHLAllRoster.json")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var teamRosters []nhl.TeamRosterDB
	if err := json.Unmarshal(fileData, &teamRosters); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	for _, roster := range teamRosters {
		if err := l.storage.RosterLoaderToDB(roster); err != nil {
			return fmt.Errorf("failed to load roster for team %s: %w", roster.TeamName, err)
		}
	}

	return nil
}

func (l *NHLLoadService) ScheduleLoader() error {
	scheduleData, err := readNHLScheduleFromJSON("./jsondata/nhl/NHLSchedule_9_04_2025.json")
	if err != nil {
		return fmt.Errorf("failed to read schedule: %w", err)
	}
	var schedules []nhl.Schedule
	for _, jsonItem := range scheduleData {
		schedule := nhl.ParseScheduleFromJSON(jsonItem)
		schedules = append(schedules, schedule)
	}

	err = l.storage.ScheduleLoaderToDB(schedules)
	if err != nil {
		return fmt.Errorf("failed to load schedule: %w", err)
	}

	return nil
}

func (l *NHLLoadService) ExportScheduleToExcel(filePath string) error {
	const op = "internal/service/nhlLoad.go"

	schedules, err := l.storage.GetAllSchedule()
	if err != nil {
		return fmt.Errorf("failed to get all schedule: %w \n in %s", err, op)
	}

	nf := excelize.NewFile()
	defer func() {
		if err := nf.Close(); err != nil {
			l.log.Error("failed to close file", slog.String("error", err.Error()))
		}
	}()

	sheetName := "NHL Schedule"
	index, err := nf.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w \n in %s", err, op)
	}

	headers := []string{
		"ID",
		"Date",
		"Time",
		"Visitor Team",
		"Home Team",
		"Visitor Score",
		"Home Score",
		"Attendance",
		"Game Duration",
		"Is Overtime?",
		"Venue",
		"Notes",
		"Created At",
		"Updated At",
	}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i) // Convert column index to Excel column name
		nf.SetCellValue(sheetName, cell, header)
	}

	style, err := nf.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDEBF7"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create style: %w \n in %s", err, op)
	}

	nf.SetCellStyle(sheetName, "A1", string(rune('A'+len(headers)-1)), style)

	for i, schedule := range schedules {
		row := i + 2 // Start from row 2 (after headers)

		nf.SetCellValue(sheetName, fmt.Sprintf("A%d", row), schedule.ID)

		dateObj, err := time.Parse(time.RFC3339, schedule.Date)
		if err == nil {
			nf.SetCellValue(sheetName, fmt.Sprintf("B%d", row), dateObj.Format("2006-01-02"))
		} else {
			nf.SetCellValue(sheetName, fmt.Sprintf("B%d", row), schedule.Date)
		}

		timeObj, err := time.Parse(time.RFC3339, schedule.Time)
		if err == nil {
			nf.SetCellValue(sheetName, fmt.Sprintf("C%d", row), timeObj.Format("15:04"))
		} else {
			nf.SetCellValue(sheetName, fmt.Sprintf("C%d", row), schedule.Time)
		}
		nf.SetCellValue(sheetName, fmt.Sprintf("D%d", row), schedule.VisitorTeam)
		nf.SetCellValue(sheetName, fmt.Sprintf("E%d", row), schedule.HomeTeam)

		// Null values handling
		var visitorScore, homeScore, attendance interface{}
		if schedule.VisitorScore != nil {
			visitorScore = *schedule.VisitorScore
		} else {
			visitorScore = ""
		}
		if schedule.HomeScore != nil {
			homeScore = *schedule.HomeScore
		} else {
			homeScore = ""
		}
		if schedule.Attendance != nil {
			attendance = *schedule.Attendance
		} else {
			attendance = ""
		}

		nf.SetCellValue(sheetName, fmt.Sprintf("F%d", row), visitorScore) // Was schedule.VisitorScore
		nf.SetCellValue(sheetName, fmt.Sprintf("G%d", row), homeScore)
		nf.SetCellValue(sheetName, fmt.Sprintf("H%d", row), attendance)

		var gameDuration, venue, notes string
		if schedule.GameDuration != nil {
			gameDuration = *schedule.GameDuration
		}
		if schedule.Venue != nil {
			venue = *schedule.Venue
		}
		if schedule.Notes != nil {
			notes = *schedule.Notes
		}

		nf.SetCellValue(sheetName, fmt.Sprintf("I%d", row), gameDuration)
		nf.SetCellValue(sheetName, fmt.Sprintf("J%d", row), schedule.IsOvertime)

		nf.SetCellValue(sheetName, fmt.Sprintf("K%d", row), venue)
		nf.SetCellValue(sheetName, fmt.Sprintf("L%d", row), notes)

		nf.SetCellValue(sheetName, fmt.Sprintf("M%d", row), schedule.CreatedAt.Format("2006-01-02 15:04:05"))
		nf.SetCellValue(sheetName, fmt.Sprintf("N%d", row), schedule.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		width := 15.0 // Init default width

		// Specific adjustments for specific columns
		switch i {
		case 3, 4: // Team names
			width = 20.0
		case 11: // Notes
			width = 30.0
		case 12, 13: // Timestamps
			width = 20.0
		}

		nf.SetColWidth(sheetName, col, col, width)
	}

	nf.SetActiveSheet(index)

	nf.DeleteSheet("Sheet1")

	if err := nf.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %w \n in %s", err, op)
	}

	l.log.Info("Schedule exported to Excel successfully", slog.String("file_path", filePath))
	return nil
}

func (l *NHLLoadService) ImportScheduleFromExcel(filePath string) error {
	const op = "internal/service/nhlLoad.go"

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w \n in %s", err, op)
	}
	defer f.Close()

	rows, err := f.GetRows("NHL Schedule")
	if err != nil {
		return fmt.Errorf("failed to get rows from sheet: %w \n in %s", err, op)
	}

	if len(rows) < 2 {
		return fmt.Errorf("no data rows in Excel file")
	}

	var schedules []nhl.ScheduleImport
	for i, row := range rows[1:] {
		if len(row) < 12 {
			continue
		}

		date, err := time.Parse("2006-01-02", row[1])
		if err != nil {
			return fmt.Errorf("invalid date format in row %d: %w \n in %s", i+2, err, op)
		}

		timeObj, err := time.Parse("15:04", row[2])
		if err != nil {
			return fmt.Errorf("invalid time format in row %d: %w \n in %s", i+2, err, op)
		}

		schedule := nhl.ScheduleImport{
			Date:         date.Format("2006-01-02"),
			Time:         timeObj.Format("15:04"),
			VisitorTeam:  row[3],
			HomeTeam:     row[4],
			VisitorScore: parseNullableInt(row[5]),
			HomeScore:    parseNullableInt(row[6]),
			Attendance:   parseNullableInt(row[7]),
			GameDuration: parseNullableString(row[8]),
			IsOvertime:   strings.ToLower(row[9]) == "true",
			Venue:        parseNullableString(row[10]),
			Notes:        parseNullableString(row[11]),
		}

		schedules = append(schedules, schedule)
	}

	if err := l.storage.UpsertSchedule(schedules); err != nil {
		return fmt.Errorf("failed to upsert schedule: %w \n in %s", err, op)
	}

	l.log.Info("Schedule imported from Excel successfully",
		slog.String("file_path", filePath),
		slog.Int("records_processed", len(schedules)))
	return nil
}

func readNHLAbbrFromJSON(filepath string) ([]nhl.TeamDB, error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var teams []nhl.TeamDB
	if err := json.Unmarshal(fileData, &teams); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return teams, nil
}

func readNHLScheduleFromJSON(filepath string) ([]nhl.ScheduleJSON, error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var schedule []nhl.ScheduleJSON
	if err := json.Unmarshal(fileData, &schedule); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return schedule, nil
}

func parseNullableInt(s string) *int {
	if s == "" {
		return nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &val
}

func parseNullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
