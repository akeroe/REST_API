package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shorneter/internal/config"
	"url-shorneter/internal/config/lib/logger/sl"
	"url-shorneter/internal/config/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter() // todo  init router chi

	// middleware - цепочка прочих хендлеров
	router.Use(middleware.RequestID) // middleware.RequestID - генерирует уникальный идентификатор для каждого запроса
	router.Use(middleware.Logger)    // middleware.Logger - логирует информацию о каждом запросе
	router.Use(mwLogger.New(log))    // mwLogger.New(log) - логирует информацию о каждом запросе с использованием slog
	router.Use(middleware.Recoverer) // middleware.Recoverer - восстанавливает приложение после паники и возвращает 500 Internal Server Error
	router.Use(middleware.URLFormat)  // middleware.URLFormat - позволяет использовать форматирование URL, например, /path.json или /path.xml

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
