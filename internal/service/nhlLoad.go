package service

import (
	"AdminAppForDiplom/internal/domain/nhlFiles"
	"AdminAppForDiplom/internal/models/nhl"
	"AdminAppForDiplom/internal/storage"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"os"
	"regexp"
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

// AbbrLoader loads NHL teams abbreviations from JSON file and inserts them into the database.
func (l *NHLLoadService) AbbrLoader() error {
	teams, err := readNHLAbbrFromJSON(nhlFiles.NHLAbbrJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read abbr: %w", err)
	}

	return l.storage.NHLLoader(teams)
}

// RosterLoader loads NHL team rosters from JSON file and inserts them into the database.
func (l *NHLLoadService) RosterLoader() error {
	teamRosters, err := readNHLAllRosterFromJSON(nhlFiles.NHLAllRosterJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read roster: %w", err)
	}

	for _, roster := range teamRosters {
		if err := l.storage.RosterLoaderToDB(roster); err != nil {
			return fmt.Errorf("failed to load roster for team %s: %w", roster.TeamName, err)
		}
	}

	return nil
}

// ScheduleLoader loads NHL schedule from JSON file and inserts it into the database.
func (l *NHLLoadService) ScheduleLoader(filePath string) error {
	scheduleData, err := readNHLScheduleFromJSON(filePath) //"./jsondata/nhl/NHLSchedule_9_04_2025.json"
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

// ExportScheduleToExcel exports NHL schedule from DB to Excel file.
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

		visitorScore := formatNullableIntForExcel(schedule.VisitorScore)
		homeScore := formatNullableIntForExcel(schedule.HomeScore)
		attendance := formatNullableIntForExcel(schedule.Attendance)

		gameDuration := formatNullableString(schedule.GameDuration)
		venue := formatNullableString(schedule.Venue)
		notes := formatNullableString(schedule.Notes)

		nf.SetCellValue(sheetName, fmt.Sprintf("A%d", row), schedule.ID)
		nf.SetCellValue(sheetName, fmt.Sprintf("D%d", row), schedule.VisitorTeam)
		nf.SetCellValue(sheetName, fmt.Sprintf("E%d", row), schedule.HomeTeam)

		nf.SetCellValue(sheetName, fmt.Sprintf("F%d", row), visitorScore) // Was schedule.VisitorScore
		nf.SetCellValue(sheetName, fmt.Sprintf("G%d", row), homeScore)
		nf.SetCellValue(sheetName, fmt.Sprintf("H%d", row), attendance)
		nf.SetCellValue(sheetName, fmt.Sprintf("I%d", row), gameDuration)

		nf.SetCellValue(sheetName, fmt.Sprintf("J%d", row), schedule.IsOvertime)

		nf.SetCellValue(sheetName, fmt.Sprintf("K%d", row), venue)
		nf.SetCellValue(sheetName, fmt.Sprintf("L%d", row), notes)

		nf.SetCellValue(sheetName, fmt.Sprintf("M%d", row), schedule.CreatedAt.Format("2006-01-02 15:04:05"))
		nf.SetCellValue(sheetName, fmt.Sprintf("N%d", row), schedule.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		width := 15.0

		// CUSTOM width for specific columns
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

// ImportScheduleFromExcel imports NHL schedule from Excel file to DB.
// Main purpose is to update schedule from Excel file to DB.
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

		date, err := parseExcelDate(row[1])
		if err != nil {
			return fmt.Errorf("invalid date format in row %d: %w \n in %s", i+2, err, op)
		}

		timeObj, err := parseExcelTime(row[2])
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

// AddNewMatchDataFromExcel adds new match data from Excel file to DB.
// Makes new Excel tempFile with new match data.
func (l *NHLLoadService) AddNewMatchDataFromExcel() error {
	const op = "internal/service/nhlLoad.AddNewMatchDataFromExcel"

	outputPath := "./temp"
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w \n in %s", err, op)
	}

	filePath := fmt.Sprintf("%s/new_nhl_match_%s.xlsx", outputPath, time.Now().Format("20060102_150405"))

	f := excelize.NewFile()
	sheetName := "New NHL Games"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w \n in %s", err, op)
	}

	headers := []string{
		"Date (YYYY-MM-DD)",
		"Time (HH:MM)",
		"Visitor Team",
		"Home Team",
		"Visitor Score",
		"Home Score",
		"Attendance",
		"Game Duration",
		"Is Overtime (true/false)",
		"Venue",
		"Notes",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i) // A1, B1, C1 и т.д.
		f.SetCellValue(sheetName, cell, header)

		switch i {
		case 0: // Date
			if err := f.AddComment(sheetName, excelize.Comment{
				Cell:   cell,
				Author: "System",
				Text:   "Формат: YYYY-MM-DD\nПример: 2024-05-15",
			}); err != nil {
				l.log.Warn("Failed to add comment for date cell", slog.String("error", err.Error()))
			}
		case 1: // Time
			if err := f.AddComment(sheetName, excelize.Comment{
				Cell:   cell,
				Author: "System",
				Text:   "Формат: HH:MM или HH:MM AM/PM\nПримеры: 19:30 или 7:30 PM",
			}); err != nil {
				l.log.Warn("Failed to add comment for time cell", slog.String("error", err.Error()))
			}
		case 8: // OT
			if err := f.AddComment(sheetName, excelize.Comment{
				Cell:   cell,
				Author: "System",
				Text:   "Введите 'true' если был овертайм\nили 'false' если не было",
			}); err != nil {
				l.log.Warn("Failed to add comment for overtime cell", slog.String("error", err.Error()))
			}
		}
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#DDEBF7"}, Pattern: 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create header style: %w \n in %s", err, op)
	}
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), headerStyle)

	exampleStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFEB9C"}, Pattern: 1},
	})
	if err != nil {
		return fmt.Errorf("failed to create example style: %w \n in %s", err, op)
	}

	colWidths := map[string]float64{
		"A": 15, "B": 12, "C": 20, "D": 20, "E": 12,
		"F": 12, "G": 12, "H": 15, "I": 18, "J": 15, "K": 30,
	}
	for col, width := range colWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	exampleRow := []interface{}{
		"2024-05-15",                  // Дата
		"19:30",                       // Время
		"Boston Bruins",               // Гостевая команда
		"Toronto Maple Leafs",         // Домашняя команда
		3,                             // Голы гостей
		2,                             // Голы дома
		18500,                         // Посещаемость
		"2:30",                        // Длительность матча
		"true",                        // Овертайм
		"Scotiabank Arena",            // Арена
		"EXAMPLE ROW - DO NOT IMPORT", // Примечания
	}
	for i, value := range exampleRow {
		cell := fmt.Sprintf("%c2", 'A'+i)
		f.SetCellValue(sheetName, cell, value)
	}

	f.SetCellStyle(sheetName, "A2", fmt.Sprintf("%c2", 'A'+len(headers)-1), exampleStyle)

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %w \n in %s", err, op)
	}

	l.log.Info("New match template created", slog.String("file_path", filePath))
	fmt.Println("\nTemplate file created at:", filePath)
	fmt.Println("Please:")
	fmt.Println("1. Open the file and add your match data below the yellow example row")
	fmt.Println("2. Save the file")
	fmt.Println("3. Return here and press Enter to import")
	fmt.Println("Or press 'q' to quit without importing")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Please don't forget to save file!\n" + "> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "":
			schedules, err := readNewGameFromExcel(filePath, sheetName)
			if err != nil {
				if err := os.Remove(filePath); err != nil {
					l.log.Warn("failed to remove temp file", slog.String("error", err.Error()))
				}
				return fmt.Errorf("failed to read match data: %w \n in %s", err, op)
			}

			if len(schedules) == 0 {
				fmt.Println("No valid match data found. Please check the file and try again.")
				continue
			}

			if err := l.storage.AddNewSchedule(schedules); err != nil {
				return fmt.Errorf("failed to add new schedule: %w \n in %s", err, op)
			}

			if err := os.Remove(filePath); err != nil {
				l.log.Warn("failed to remove temp file", slog.String("error", err.Error()))
			}

			fmt.Println("New match(es) added successfully!")
			return nil

		case "q", "quit":
			fmt.Println("Operation cancelled. Template file remains at:", filePath)
			if err := os.Remove(filePath); err != nil {
				l.log.Warn("failed to remove temp file", slog.String("error", err.Error()))
			}
			return nil

		default:
			fmt.Println("Invalid input. Press 'Enter' to continue or 'q' to quit")
		}
	}
}

// readNHLAbbrFromJSON reads a new abbreviation from an Excel file.
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

// readNHLAllRosterFromJSON reads a new roster from an Excel file.
func readNHLAllRosterFromJSON(filepath string) ([]nhl.TeamRosterDB, error) {
	const op = "internal/service/nhlLoad.go/readNHLAllRosterFromJSON"

	fileDate, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w \n in %s", err, op)
	}

	var teamRosters []nhl.TeamRosterDB
	if err := json.Unmarshal(fileDate, &teamRosters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w \n in %s", err, op)
	}

	return teamRosters, nil
}

// readNHLScheduleFromJSON reads a new schedule from an Excel file.
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

func formatNullableIntForExcel(value *int) interface{} {
	if value != nil {
		return *value
	}
	return ""
}

func formatNullableString(value *string) string {
	if value != nil {
		return *value
	}
	return ""
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

// parseExcelDate try to parse date from string
// it will be made by AI and will be used in the future
// I'm sorry for this function, but I don't know how to parse date from string
func parseExcelDate(dateStr string) (time.Time, error) {
	// Trim any whitespace
	dateStr = strings.TrimSpace(dateStr)

	// Пробуем разные форматы даты
	formats := []string{
		"2006-01-02",     // ISO формат (2024-10-04)
		"02.01.2006",     // Европейский формат (04.10.2024)
		"01/02/2006",     // Американский формат (10/04/2024)
		"2 Jan 2006",     // Текстовый формат с сокращенным месяцем
		"2 January 2006", // Текстовый формат с полным месяцем
		"02 01 2006",     // Формат с пробелами (04 10 2024)
		"02 01 06",       // Сокращенный формат с пробелами (04 10 24)
		"01-02-06",       // MM-DD-YY формат (05-15-25)
		"01-02-2006",     // MM-DD-YYYY формат (05-15-2025)
		"2006-01-02",     // YYYY-MM-DD формат (2025-05-15)
		"02-01-2006",     // DD-MM-YYYY формат (15-05-2025)
		"02-01-06",       // DD-MM-YY формат (15-05-25)
	}

	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			if format == "01-02-06" || format == "02-01-06" || format == "02 01 06" {
				if date.Year() < 2000 {
					date = date.AddDate(100, 0, 0) // Add 100 years for 2-digit years before 2000
				}
			}
			return date, nil
		}
	}

	if matched, _ := regexp.MatchString(`^\d{2}-\d{2}-\d{2}$`, dateStr); matched {
		// For XX-XX-XX format, try to intelligently determine if it's MM-DD-YY or DD-MM-YY
		parts := strings.Split(dateStr, "-")
		if len(parts) == 3 {
			month, _ := strconv.Atoi(parts[0])
			day, _ := strconv.Atoi(parts[1])
			year, _ := strconv.Atoi(parts[2])

			// Add 2000 to convert 2-digit year to 4-digit
			if year < 100 {
				year += 2000
			}

			// If month value is valid (1-12), assume MM-DD-YY
			if month >= 1 && month <= 12 && day >= 1 && day <= 31 {
				return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
			}
		}
	}

	// Если не удалось распознать дату, пробуем извлечь числовое значение Excel
	if numVal, err := strconv.ParseFloat(dateStr, 64); err == nil {
		// Преобразуем числовое значение Excel в дату
		// Excel начинает отсчет с 1 января 1900, но есть ошибка с високосным годом
		// поэтому вычитаем 2 дня для корректировки
		excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		days := int(numVal)
		return excelEpoch.AddDate(0, 0, days), nil
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// parseExcelTime
func parseExcelTime(timeStr string) (time.Time, error) {
	timeFormats := []string{
		"15:04",    // 24 (14:30)
		"3:04 PM",  // 12 with AM/PM (2:30 PM)
		"3:04PM",   // 12 without (_) AM/PM (2:30PM)
		"15:04:05", // 24 with seconds (14:30:05)
	}

	for _, format := range timeFormats {
		timeObj, err := time.Parse(format, timeStr)
		if err == nil {
			return timeObj, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// readNewGameFromExcel in this function we read data from excel file and return slice of new games
// All of these data will be saved in database
func readNewGameFromExcel(filePath, sheetName string) ([]nhl.ScheduleImport, error) {
	const op = "internal/service/nhlLoad.go.readNewGameFromExcel"

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w \n in %s", err, op)
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet: %w \n in %s", err, op)
	}

	if len(rows) < 3 { // Need at least header + example + one data row
		return nil, fmt.Errorf("no data rows in Excel file")
	}

	var newGames []nhl.ScheduleImport
	for i, row := range rows[2:] { // Skip header and example rows
		if len(row) < 11 {
			continue
		}

		if len(row) > 10 && strings.Contains(row[10], "EXAMPLE ROW") {
			continue
		}

		date, err := parseExcelDate(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid date format in row %d: %w \n in %s", i+3, err, op)
		}

		timeObj, err := parseExcelTime(row[1])
		if err != nil {
			return nil, fmt.Errorf("invalid time format in row %d: %w \n in %s", i+3, err, op)
		}

		schedule := nhl.ScheduleImport{
			Date:         date.Format("2006-01-02"),
			Time:         timeObj.Format("15:04"),
			VisitorTeam:  row[2],
			HomeTeam:     row[3],
			VisitorScore: parseNullableInt(row[4]),
			HomeScore:    parseNullableInt(row[5]),
			Attendance:   parseNullableInt(row[6]),
			GameDuration: parseNullableString(row[7]),
			IsOvertime:   strings.ToLower(row[8]) == "true",
			Venue:        parseNullableString(row[9]),
			Notes:        parseNullableString(row[10]),
		}

		fmt.Println(schedule)

		newGames = append(newGames, schedule)
	}

	return newGames, nil
}
