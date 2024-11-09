package main

import (
	"AdminAppForDiplom/cmd/cobra"
	"AdminAppForDiplom/pkg/lib/logger/handler/slogpretty"
	"log/slog"
	"os"
)

func main() {
	logg := setupPrettySlogLocal()

	logg.Info("app starting")

	cobra.Execute(logg)
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
