package cobra

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

var rootCmd = &cobra.Command{
	Short: "Admin app for SportThunder app",
}

func Execute(log *slog.Logger) {
	if err := rootCmd.Execute(); err != nil {
		handleErr(log, err)
	}
}

func handleErr(log *slog.Logger, err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	log.Error("failed to execute command", slog.Any("error", err))
	os.Exit(1)
}
