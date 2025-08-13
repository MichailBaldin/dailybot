package main

import (
	"dailybot/internal/bot"
	"dailybot/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	b, err := bot.New(cfg)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	b.Start()
}
