package cobra

import (
	"AdminAppForDiplom/internal/domain/direct"
	"AdminAppForDiplom/internal/parser"
	"AdminAppForDiplom/internal/parser/nba"
	"AdminAppForDiplom/internal/parser/nhl"
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/spf13/cobra"
)

var nbaCmd = &cobra.Command{
	Use:   "nba",
	Short: "NBA",
	Long:  "Main NBA command which contains subcommands",
}

var nbaParseCmd = &cobra.Command{
	Use:   "nba-prs",
	Short: "NBA parser",
	Long:  "It is NBA parser",
	Run:   parseNBA,
}

var nbaLoadToDBCmd = &cobra.Command{
	Use:   "nba-db",
	Short: "NBA loader",
	Long:  "It is NBA loader to DB",
	Run:   loadNBAToDB,
}

func parseNBA(cmd *cobra.Command, _ []string) {
	ctx.log.Info("NBA parse started")

	var objParse func(g *geziyor.Geziyor, r *client.Response)
	var filePath string
	var directObj string

	method := cmd.Flag("method").Value.String()
	fmt.Printf("method: %s\n", method)
	ctx.log.Info("method", "method", method)

	switch method {
	case "abbr":
		//TODO: add to config
		filePath = jsonPath + nbaPath + abbrNBA
		objParse = nba.NameParse
		directObj = direct.NBANameAbbr

	case "roster":
		filePath = jsonPath + nbaPath + rosterNBA
		objParse = nhl.RosterParse
		directObj = direct.NHLRoster

	case "allroster":
		filePath = jsonPath + nbaPath + allRosterNBA
		objParse = nhl.AllRosterParse
		directObj = direct.AllNHLRoster

	case "debug":
		filePath = jsonPath + nbaPath + debugRosterNBA
		//objParse = nhl.AllRosterParse
		//directObj = direct.AllNHLRoster

	default:
		ctx.log.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	parser.Parser(objParse, filePath, directObj)

	ctx.log.Info("NBA parse finished")
}

func loadNBAToDB(cmd *cobra.Command, _ []string) {
	ctx.log.Info("NBA loader started")

	method := cmd.Flag("method").Value.String()
	ctx.log.Info("method", "method", method)

	switch method {
	case "abbr":
		ctx.log.Info("NBA loader started")
		err := ctx.service.NBALoad.AbbrNBALoader()
		if err != nil {
			ctx.log.Error("failed to load abbreviationNHL to DB")
		}
	default:
		ctx.log.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	ctx.log.Info("NBA loader finished")
}
