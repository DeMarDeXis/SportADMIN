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
