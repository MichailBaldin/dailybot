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
	AdminPort      string
	AdminPassword  string
}

func Load() (*Config, error) {
	// Пытаемся загрузить .env файл (для локальной разработки)
	// В продакшене этот файл может отсутствовать - это нормально
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Info: .env file not found (using environment variables): %v\n", err)
	}

	cfg := &Config{
		TelegramToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		OpenWeatherKey: os.Getenv("OPENWEATHER_API_KEY"),
		NewsAPIKey:     os.Getenv("NEWS_API_KEY"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AdminPort:      getEnvWithDefault("ADMIN_PORT", "8080"),
		AdminPassword:  getEnvWithDefault("ADMIN_PASSWORD", "admin123"),
	}

	// Проверяем обязательные переменные
	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	// Устанавливаем значения по умолчанию
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://dailybot:password@localhost:5432/dailybot?sslmode=disable"
	}

	// Логируем какие API ключи настроены (без вывода самих ключей)
	fmt.Printf("Config loaded:\n")
	fmt.Printf("- Telegram Bot: configured\n")
	fmt.Printf("- OpenWeather API: %s\n", getStatus(cfg.OpenWeatherKey))
	fmt.Printf("- News API: %s\n", getStatus(cfg.NewsAPIKey))
	fmt.Printf("- Admin Panel: port %s\n", cfg.AdminPort)

	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getStatus(value string) string {
	if value == "" {
		return "not configured (demo mode)"
	}
	return "configured"
}
