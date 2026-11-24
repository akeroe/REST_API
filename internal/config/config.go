package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// Получаем путь до конфиг-файла из env-переменной CONFIG_PATH
func MustLoad() *Config {
	err := godotenv.Load() // читает .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Проверяем существование конфиг-файла
	configPath := os.Getenv("CONFIG_PATH") // os.Getenv берет из .env CONFIG_PATH
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	var cfg Config
	// Читаем конфиг-файл и заполняем нашу структуру
	err = cleanenv.ReadConfig(configPath, &cfg) // err != nil - проверить: ошибка есть ? и если есть, то фатал
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}
