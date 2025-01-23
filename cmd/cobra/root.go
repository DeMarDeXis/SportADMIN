package cobra

import (
	"AdminAppForDiplom/internal/service"
	storage2 "AdminAppForDiplom/internal/storage"
	"AdminAppForDiplom/internal/storage/postgres"
	"AdminAppForDiplom/pkg/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

type commandContext struct {
	db      *sqlx.DB
	log     *slog.Logger
	config  *config.Config
	service *service.Service
}

var ctx = &commandContext{}

var rootCmd = &cobra.Command{
	Short: "Admin app for SportThunder app",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		db, err := postgres.New(ctx.config.DB, ctx.log)
		if err != nil {
			handleErr(ctx.log, fmt.Errorf("failed to init database: %w", err))
		}
		ctx.db = db

		storage := storage2.NewStorage(db)
		ctx.service = service.NewService(storage, ctx.log)
	},
}

func Execute(log *slog.Logger, cfg *config.Config) {
	ctx.log = log
	ctx.config = cfg
	if err := rootCmd.Execute(); err != nil {
		handleErr(log, err)
	}
}

func handleErr(log *slog.Logger, err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	log.Error("failed to execute command", slog.Any("error", err))
	os.Exit(1)
}
