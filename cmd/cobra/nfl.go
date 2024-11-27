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

func parseNFL(cmd *cobra.Command, _ []string) {
	slog.Info("NHL parse started")

	var objParse func(g *geziyor.Geziyor, r *client.Response)
	var filePath string
	var directObj string

	method := cmd.Flag("method").Value.String()
	fmt.Printf("method: %s\n", method)
	slog.Info("method", "method", method)

	switch method {
	case "nfl":
		filePath = jsonPath + "NFLAbbr.json"
		objParse = nfl.ParseAbbr
		directObj = direct.NFLNameAbbr
	default:
		slog.Error("Unsupported method", "method", method)
		cmd.PrintErr("Unknown method")
	}

	parser.Parser(objParse, filePath, directObj)

	slog.Info("NHL parse finished")
}
