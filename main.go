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
// TODO: think about the cases of the absolute and relative paths(Import cmd) (multiple cases) (15.04.25)

// TODO: deleter in postgres (08.05.25)

// TODO: add GoDoc comments to all functions in storage(11.05.25)

// TODO: add MLB teams to DB (11.05.25)

//TODO: fix wrapping errores (12.06.25) by example nhl/roster cmd
