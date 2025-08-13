package bot

import (
	"dailybot/internal/config"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	config *config.Config
}

func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, err
	}

	api.Debug = false
	log.Printf("Bot authorized as @%s", api.Self.UserName)

	return &Bot{
		api:    api,
		config: cfg,
	}, nil
}

func (b *Bot) Start() {
	log.Println("Bot started and listening for updates...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID

	switch message.Command() {
	case "start":
		b.handleStart(chatID)
	case "help":
		b.handleHelp(chatID)
	case "weather":
		b.handleWeather(chatID, message.CommandArguments())
	case "news":
		b.handleNews(chatID)
	default:
		if message.IsCommand() {
			b.sendMessage(chatID, "Unknown command. Please use /help")
		}
	}
}

func (b *Bot) sendMessage(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "HTML"
	b.api.Send(msg)
}
