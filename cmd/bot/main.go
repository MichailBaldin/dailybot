package main

import (
	"dailybot/internal/admin"
	"dailybot/internal/bot"
	"dailybot/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Создаем простую админку
	adminServer := admin.NewSimpleAdmin(cfg)

	// Создаем бота
	b, err := bot.New(cfg, adminServer)
	if err != nil {
		log.Fatal("Failed to create bot:", err)
	}

	// Запускаем админку в отдельной горутине
	go func() {
		log.Println("Starting admin panel...")
		adminServer.Start()
	}()

	// Запускаем бота (блокирующий вызов)
	log.Println("Starting Telegram bot...")
	b.Start()
}
