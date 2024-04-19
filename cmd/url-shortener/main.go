package main

import (
	"os"

	"github.com/Svoevolin/url-shortener/internal/config"
	"github.com/Svoevolin/url-shortener/internal/database/postgres"
	"github.com/Svoevolin/url-shortener/internal/database/postgres/models"
	"github.com/Svoevolin/url-shortener/internal/lib/logger/sl"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	// TO DO: init database: postgresql

	db, err := postgres.New(cfg.Dsn)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	urlDB := models.NewUrlDB(db)
	_ = urlDB

	// TO DO: init router: chi, "chi render"

	// TO DO: run server
}

func setupLogger(env string) (log *slog.Logger) {

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}

	return log
}
