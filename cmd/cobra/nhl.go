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

// TODO: add to config
const (
	jsonPath = "./jsondata/"
)

func parseNHL(cmd *cobra.Command, _ []string) {
	slog.Info("NHL parse started")

	var objParse func(g *geziyor.Geziyor, r *client.Response)
	var filePath string
	var directObj string

	method := cmd.Flag("method").Value.String()
	fmt.Printf("method: %s\n", method)
	slog.Info("method", "method", method)

	switch method {
	case "abbr":
		//TODO: add to config
		filePath = jsonPath + "NHLAbbr.json"
		objParse = nhl.NameParse
		directObj = direct.NHLNameAbbr

	case "roster":
		filePath = jsonPath + "NHLRoster.json"
		objParse = nhl.RosterParse
		directObj = direct.NHLRoster

	case "allroster":
		filePath = jsonPath + "NHLAllRoster.json"
		objParse = nhl.AllRosterParse
		directObj = direct.AllNHLRoster
	default:
		slog.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	parser.Parser(objParse, filePath, directObj)

	slog.Info("NHL parse finished")
}

func loadNHLToDB(cmd *cobra.Command, _ []string) {
	ctx.log.Info("NHL loader started")

	method := cmd.Flag("method").Value.String()
	ctx.log.Info("method", "method", method)

	switch method {
	case "abbr":
		err := ctx.service.Load.Loader()
		if err != nil {
			ctx.log.Error("failed to load to DB", "error", err)
			return
		}

	case "roster":

	case "allroster":

	default:
		slog.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	ctx.log.Info("NHL loader finished")
}
