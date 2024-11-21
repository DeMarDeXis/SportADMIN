package cobra

import (
	"AdminAppForDiplom/internal/parser/nhl"
	"AdminAppForDiplom/pkg/lib/customjsonexp"
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"github.com/spf13/cobra"
	"log/slog"
)

const (
	nhlNameAbbr = "https://en.wikipedia.org/wiki/Wikipedia:WikiProject_Ice_Hockey/NHL_team_abbreviations"
)

const (
	jsonPath = "./jsondata/"
)

var nhlCmd = &cobra.Command{
	Use:   "nhl",
	Short: "NHL",
	Run:   parse,
}

func parse(cmd *cobra.Command, _ []string) {
	slog.Info("NHL parse started")

	var filePath string
	var parser func(g *geziyor.Geziyor, r *client.Response)

	method := cmd.Flag("method").Value.String()

	switch method {
	case "abbr":
		filePath = jsonPath + "NHLAbbr.json"
		parser = nhl.NHLNameParse

	case "roster":
		//filePath = jsonPath + "NHLAbbr.json"
		fmt.Println("roster")
		return

	default:
		slog.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	exporter, err := customjsonexp.NewCustomJSONExporter(filePath)
	if err != nil {
		slog.Error("Error creating exporter", "error", err)
		return
	}

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{nhlNameAbbr},
		ParseFunc: parser,
		Exporters: []export.Exporter{exporter},
	}).Start()

	slog.Info("NHL parse finished")
}

func init() {
	rootCmd.AddCommand(nhlCmd)

	nhlCmd.Flags().StringP("method", "m", "", "What are we parsing")

	if err := nhlCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
