package main

import (
	"AdminAppForDiplom/cmd/cobra"
	"AdminAppForDiplom/pkg/config"
	"AdminAppForDiplom/pkg/lib/logger/handler/slogpretty"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

func main() {
	logg := setupPrettySlogLocal()
	logg.Debug("Logg initialized")

	cfg, err := config.LoadConfig("")
	if err != nil {
		logg.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}
	logg.Debug("config loaded", slog.Any("config", cfg))

	logg.Info("app starting")
	cobra.Execute(logg, cfg)
}

func setupPrettySlogLocal() *slog.Logger {
	opts := slogpretty.PrettyHandlersOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handlerLog := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handlerLog)
}

// TODO: fix t.j. oshie err
// TODO: add flags-methods to drawio
// TODO: add img_url to teams and push to github
