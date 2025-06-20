package parser

import (
	"AdminAppForDiplom/pkg/lib/customjsonexp"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"log/slog"
)

func Parser(parse func(g *geziyor.Geziyor, r *client.Response), filePath string, direct string) {
	exporter, err := customjsonexp.NewCustomJSONExporter(filePath)
	if err != nil {
		slog.Error("Error creating exporter", "error", err)
		return
	}

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{direct},
		ParseFunc: parse,
		Exporters: []export.Exporter{exporter},
	}).Start()
}
