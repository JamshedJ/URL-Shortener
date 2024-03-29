package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/JamshedJ/URL-Shortener/internal/config"
	"github.com/JamshedJ/URL-Shortener/internal/http-server/handlers/fetch"
	"github.com/JamshedJ/URL-Shortener/internal/http-server/handlers/save"
	mwLogger "github.com/JamshedJ/URL-Shortener/internal/http-server/middleware/logger"
	"github.com/JamshedJ/URL-Shortener/internal/lib/logger/sl"
	"github.com/JamshedJ/URL-Shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Print something")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID) // middleware из chi, каждому приходящему запросу присваивает RequestID, полезен для отслеживание трейсов
	router.Use(mwLogger.New(log))    // middleware для логгирование запросов
	router.Use(middleware.Recoverer) // ловит панику

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", fetch.New(log, storage))
	
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error runnig server")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
