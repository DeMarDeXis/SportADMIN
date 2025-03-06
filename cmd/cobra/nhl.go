package cobra

import (
	"AdminAppForDiplom/internal/domain/direct"
	"AdminAppForDiplom/internal/parser"
	"AdminAppForDiplom/internal/parser/nhl"
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/spf13/cobra"
	"log/slog"
)

func parseNHL(cmd *cobra.Command, _ []string) {
	ctx.log.Info("NHL parse started")

	var objParse func(g *geziyor.Geziyor, r *client.Response)
	var filePath string
	var directObj string

	method := cmd.Flag("method").Value.String()
	fmt.Printf("method: %s\n", method)
	ctx.log.Info("method", "method", method)

	switch method {
	case "abbr":
		filePath = jsonPath + abbrNHL
		objParse = nhl.NameParse
		directObj = direct.NHLNameAbbr

	case "roster":
		filePath = jsonPath + rosterNHL
		objParse = nhl.RosterParse
		directObj = direct.NHLRoster

	case "allroster":
		filePath = jsonPath + allRosterNHL
		objParse = nhl.AllRosterParse
		directObj = direct.AllNHLRoster

	case "debug":
		filePath = jsonPath + debugRosterNHL
		//objParse = nhl.AllRosterParse
		//directObj = direct.AllNHLRoster

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
	case "abbr":
		err := ctx.service.NHLLoad.AbbrLoader()
		if err != nil {
			ctx.log.Error("failed to load abbreviationNHL to DB",
				slog.String("error", err.Error())) // Changed this line
			return
		}
	case "allroster":
		err := ctx.service.NHLLoad.RosterLoader()
		if err != nil {
			ctx.log.Error("failed to load rosterNHL to DB", "error", err)
			return
		}

	default:
		slog.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	ctx.log.Info("NHL loader finished")
}
