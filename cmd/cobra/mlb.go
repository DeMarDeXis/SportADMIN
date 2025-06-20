package cobra

import (
	"AdminAppForDiplom/internal/domain/direct"
	"AdminAppForDiplom/internal/parser"
	"AdminAppForDiplom/internal/parser/mlb"
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/spf13/cobra"
	"log/slog"
)

var mlbCmd = &cobra.Command{
	Use:   "mlb",
	Short: "MLB",
	Long:  "Main MLB command which contains subcommands",
}

var mlbParseCmd = &cobra.Command{
	Use:   "mlb-prs",
	Short: "MLB parser",
	Long:  "It is MLB parser",
	Run:   parseMLB,
}

var mlbLoadToDBCmd = &cobra.Command{
	Use:   "mlb-db",
	Short: "MLB loader",
	Long:  "It is MLB loader to DB",
	Run:   loadMLBToDB,
}

func parseMLB(cmd *cobra.Command, _ []string) {
	ctx.log.Info("MLB parse started")

	var objParse func(g *geziyor.Geziyor, r *client.Response)
	var filePath string
	var directObj string

	method := cmd.Flag("method").Value.String()
	fmt.Printf("method: %s\n", method)
	ctx.log.Info("method", "method", method)

	switch method {
	case "abbr":
		filePath = jsonPath + mlbPath + abbrMLB
		objParse = mlb.NameParse
		directObj = direct.MLBNameAbbr

	case "roster":
		ctx.log.Info("RosterParse")

	case "allroster":
		ctx.log.Info("AllRosterParse")

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

func loadMLBToDB(cmd *cobra.Command, _ []string) {
	ctx.log.Info("MLB loader started")

	method := cmd.Flag("method").Value.String()
	ctx.log.Info("method", "method", method)

	switch method {
	case "abbr":
		err := ctx.service.MLBLoad.AbbrMLBLoader()
		if err != nil {
			ctx.log.Error("failed to load abbrMLB to DB", slog.String("error", err.Error()))
			return
		}
	default:
		ctx.log.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	ctx.log.Info("MLB loader finished")
}
