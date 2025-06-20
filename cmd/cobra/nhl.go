package cobra

import (
	"AdminAppForDiplom/internal/domain/direct"
	"AdminAppForDiplom/internal/parser"
	"AdminAppForDiplom/internal/parser/nhl"
	"bufio"
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	abbrFlagNhl           = "abbr"
	rosterFlagNhl         = "roster"
	allRosterFlagNhl      = "allroster"
	scheduleFlagNhl       = "schedule"
	exportScheduleFlagNhl = "exportScheduleNhlXls"
	updateScheduleFlagNHL = "UpdScheduleNhlXls"
	AddNewMatchFlagNHL    = "AddNewMatchNhlXls" //TODO: add this command in the Description
)

var nhlCmd = &cobra.Command{
	Use:   "nhl",
	Short: "NHL",
	Long:  "Main NHL command which contains subcommands",
}

var nhlParseCmd = &cobra.Command{
	Use:   "nhl-prs",
	Short: "NHL parser",
	Long:  "It is NHL parser",
	Run:   parseNHL,
	//TODO: add example and description like in nhlLoadToDBCmd (right down)
}

var nhlLoadToDBCmd = &cobra.Command{
	Use:   "nhl-db",
	Short: "NHL loader",
	Long: "It is NHL loader to DB" +
		"NHL loader has next subcommands: " + "\n" +
		"-" + abbrFlagNhl + "\n" +
		"-" + allRosterFlagNhl + "\n" +
		"-" + scheduleFlagNhl + "\n" +
		"-" + exportScheduleFlagNhl + "\n" +
		"-" + updateScheduleFlagNHL + "\n" +
		"-" + AddNewMatchFlagNHL + "\n",
	// TODO: add new const when new subcommand will be added
	Example: "./sportthunder nhl nhl-db -m importScheduleNhlXls",
	Run:     loadNHLToDB,
}

func parseNHL(cmd *cobra.Command, _ []string) {
	ctx.log.Info("NHL parse started")

	var objParse func(g *geziyor.Geziyor, r *client.Response)
	var filePath string
	var directObj string

	method := cmd.Flag("method").Value.String()
	fmt.Printf("method: %s\n", method)
	ctx.log.Info("method", "method", method)

	switch method {
	case abbrFlagNhl:
		filePath = jsonPath + nhlPath + abbrNHL
		objParse = nhl.NameParse
		directObj = direct.NHLNameAbbr

	case rosterFlagNhl:
		filePath = jsonPath + nhlPath + rosterNHL
		objParse = nhl.RosterParse
		directObj = direct.NHLRoster

	case allRosterFlagNhl:
		filePath = jsonPath + nhlPath + allRosterNHL
		objParse = nhl.AllRosterParse
		directObj = direct.AllNHLRoster

	default:
		ctx.log.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	parser.Parser(objParse, filePath, directObj)

	ctx.log.Info("NHL parse finished")
}

func loadNHLToDB(cmd *cobra.Command, _ []string) {
	ctx.log.Info("NHL loader started")

	method := cmd.Flag("method").Value.String()
	ctx.log.Info("method", "method", method)

	switch method {
	case abbrFlagNhl:
		err := ctx.service.NHLLoad.AbbrLoader()
		if err != nil {
			ctx.log.Error("failed to load abbreviationNHL to DB", slog.String("error", err.Error()))
			return
		}
	case allRosterFlagNhl:
		err := ctx.service.NHLLoad.RosterLoader()
		if err != nil {
			ctx.log.Error("failed to load rosterNHL to DB", "error", slog.String("error", err.Error()))
			return
		}
	case scheduleFlagNhl:
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the path to the JSON file: ")
		filePath, _ := reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)
		if !strings.HasSuffix(filePath, ".json") {
			filePath += ".json"
		}
		buildPath := "./" + filePath
		err := ctx.service.NHLLoad.ScheduleLoader(buildPath)
		if err != nil {
			ctx.log.Error("failed to load scheduleNHL to DB", slog.String("error", err.Error()))
			return
		}

	case exportScheduleFlagNhl:
		outputPath := "./exports"
		filePath := fmt.Sprintf("%s/nhl_schedule_%s.xlsx", outputPath, time.Now().Format("20060102_150405"))

		err := ctx.service.NHLLoad.ExportScheduleToExcel(filePath)
		if err != nil {
			ctx.log.Error("failed to export scheduleNHL to Excel", slog.String("error", err.Error()))
			return
		}

	case updateScheduleFlagNHL:
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the path to the Excel file: ")
		filePath, _ := reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)
		if !strings.HasSuffix(filePath, ".xlsx") {
			filePath += ".xlsx"
		}
		buildPath := "./" + filePath
		err := ctx.service.NHLLoad.ImportScheduleFromExcel(buildPath)
		if err != nil {
			ctx.log.Error("failed to import scheduleNHL to Excel", slog.String("error", err.Error()))
		}

	case AddNewMatchFlagNHL:
		err := ctx.service.AddNewMatchDataFromExcel()
		if err != nil {
			ctx.log.Error("failed to add new match", slog.String("error", err.Error()))
		}

	default:
		ctx.log.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	ctx.log.Info("NHL loader finished")
}
