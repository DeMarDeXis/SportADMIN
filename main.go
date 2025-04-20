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

// TODO: fix t.j. oshie err in nhl roster (31.03.25)
// TODO: add path from constants: internal/service/nhlLoad.go (29.03.25)
// TODO: think about the cases of the absolute and relative paths(Import cmd) (15.04.25) (partial)
// TODO: not working updated_at in ImportToPostgres(15.04.25)
