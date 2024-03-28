package main

import (
	"log/slog"
	"os"

	"github.com/JamshedJ/URL-Shortener/internal/config"
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
	// TODO: init config
	cfg := config.MustLoad()

	// TODO: init logger
	log := setupLogger(cfg.Env)
	log.Info("Print something")

	// TODO: init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage
	// TODO: init router
	router := chi.NewRouter()
 
	router.Use(middleware.RequestID) // middleware из chi, каждому приходящему запросу присваивает RequestID, полезен для отслеживание трейсов
	router.Use(mwLogger.New(log)) 	 // middleware для логгирование запросов
	// TODO: run server
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
