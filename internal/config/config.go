package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken  string
	OpenWeatherKey string
	NewsAPIKey     string
	DatabaseURL    string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		_ = fmt.Errorf("Failed to load .env")
	}

	cfg := &Config{
		TelegramToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		OpenWeatherKey: os.Getenv("OPENWEATHER_API_KEY"),
		NewsAPIKey:     os.Getenv("NEWS_API_KEY"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
	}

	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://dailybot:password@localhost:5432/dailybot?sslmode=disable"
	}

	return cfg, nil
}
