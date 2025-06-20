package cobra

import (
	"AdminAppForDiplom/internal/domain/direct"
	"AdminAppForDiplom/internal/parser"
	"AdminAppForDiplom/internal/parser/nfl"
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/spf13/cobra"
	"log/slog"
)

var nflCmd = &cobra.Command{
	Use:   "nfl",
	Short: "NFL",
}

var nflParseCmd = &cobra.Command{
	Use:   "nfl-prs",
	Short: "NFL parser",
	Long:  "It is NFL parser",
	Run:   parseNFL,
}

var nflLoadToDBCmd = &cobra.Command{
	Use:   "nfl-db",
	Short: "NFL loader",
	Long:  "It is NFL loader to DB",
	Run:   loadNFLToDB,
}


func parseNFL(cmd *cobra.Command, _ []string) {
	slog.Info("NHL parse started")

	var objParse func(g *geziyor.Geziyor, r *client.Response)
	var filePath string
	var directObj string

	method := cmd.Flag("method").Value.String()
	fmt.Printf("method: %s\n", method)
	slog.Info("method", "method", method)

	switch method {
	case "nflabbr":
		filePath = jsonPath + nflPath + abbrNFL
		objParse = nfl.ParseAbbr
		directObj = direct.NFLNameAbbr
	case "nflroster":
		slog.Info("nflroster parse started")
	default:
		slog.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	parser.Parser(objParse, filePath, directObj)

	slog.Info("NHL parse finished")
}

func loadNFLToDB(cmd *cobra.Command, _ []string) {
	ctx.log.Info("NFL loader started")

	method := cmd.Flag("method").Value.String()
	ctx.log.Info("method", "method", method)

	switch method {
	case "abbr":
		ctx.log.Info("NFL loader started")
		err := ctx.service.NFLLoad.AbbrNFLLoader()
		if err != nil {
			ctx.log.Error("failed to load abbreviationNHL to DB")
		}
	default:
		ctx.log.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	ctx.log.Info("NFL loader finished")
}
