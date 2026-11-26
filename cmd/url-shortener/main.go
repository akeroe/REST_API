package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shorneter/internal/config"
	"url-shorneter/internal/config/lib/logger/sl"
	"url-shorneter/internal/config/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad() // init config
	fmt.Println(cfg)

	log := setupLogger(cfg.Env) // todo init logger: slog
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath) // todo  init storage: sqlite -> postgres

	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1) // приложение падает с кодом ошибки
	}

	
	// todo  init router chi

	// todo run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
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
