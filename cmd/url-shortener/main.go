package main

import (
	"net/http"
	"os"

	"github.com/Svoevolin/url-shortener/internal/config"
	"github.com/Svoevolin/url-shortener/internal/database/postgres"
	"github.com/Svoevolin/url-shortener/internal/database/postgres/models"
	"github.com/Svoevolin/url-shortener/internal/http-server/handlers/url/delete"
	"github.com/Svoevolin/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/Svoevolin/url-shortener/internal/http-server/handlers/url/save"
	"github.com/Svoevolin/url-shortener/internal/http-server/handlers/url/update"
	mwLogger "github.com/Svoevolin/url-shortener/internal/http-server/middleware/logger"
	"github.com/Svoevolin/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/Svoevolin/url-shortener/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// CONFIG

	cfg := config.MustLoad()

	// LOGGER

	log := setupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env))
	log.Info("starting url-shortener", slog.String("env", cfg.Env), slog.String("version", "123"))
	log.Debug("debug messages are enabled")

	// DATABASE

	db, err := postgres.New(cfg.Dsn)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	urlDB := models.NewUrlDB(db)

	// ROUTER

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, urlDB, cfg.AliasLength))
		r.Put("/", update.New(log, urlDB, cfg.AliasLength))
		r.Delete("/{alias}", delete.New(log, urlDB))
	})

	router.Get("/{alias}", redirect.New(log, urlDB))

	log.Info("starting server", slog.String("address", cfg.Address))

	// RUN SERVER

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) (log *slog.Logger) {

	switch env {
	case envLocal:
		log = setupPrettySlog()

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
