package bot

import (
	"dailybot/internal/api"
	"fmt"
	"log"
	"strings"
)

func (b *Bot) handleStart(chatID int64) {
	text := `<b>Привет! Я ДейлиБот - твой помощник на каждый день!</b>

<b>Мои команды:</b>
/weather [город] - прогноз погоды
/exchange [валюта] - курс валют ЦБ РФ
/news - главные новости дня
/help - подробная справка`

	b.sendMessage(chatID, text)
}

func (b *Bot) handleHelp(chatID int64) {
	text := `<b>Справка по командам:</b>

<b>Мои команды:</b>
<b>/weather [город]</b> - получить прогноз погоды
Пример: <code>/weather Москва</code>

<b>/news</b> - топ-5 главных новостей дня
Актуальные новости из российских источников

<b>/exchange [валюта]</b> - курс валют по данным ЦБ РФ
Пример: <code>/exchange USD</code> или <code>/exchange EUR</code>

<i>Бот работает на языке Go и использует официальные API</i>`

	b.sendMessage(chatID, text)
}

func (b *Bot) handleWeather(chatID int64, args string) {
	city := strings.TrimSpace(args)
	if city == "" {
		b.sendMessage(chatID, "Укажите город для получения прогноза погоды\n\nПример: <code>/weather Москва</code>")
		return
	}

	b.sendMessage(chatID, "Получаю данные о погоде...")

	weatherInfo, err := api.GetWeather(city, b.config.OpenWeatherKey)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("<b>Ошибка:</b> %s", err.Error()))
		return
	}

	b.sendMessage(chatID, weatherInfo)
}

func (b *Bot) handleNews(chatID int64) {
	b.sendMessage(chatID, "Загружаю актуальные новости...")

	log.Printf("Fetching news for chat %d", chatID)

	newsInfo, err := api.GetNews(b.config.NewsAPIKey)
	if err != nil {
		log.Printf("News error: %v", err)
		b.sendMessage(chatID, fmt.Sprintf("<b>Ошибка:</b> %s", err.Error()))
		return
	}

	log.Printf("News fetched successfully")
	b.sendMessage(chatID, newsInfo)
}

func (b *Bot) handleExchange(chatID int64, args string) {
	currency := strings.TrimSpace(strings.ToUpper(args))
	if currency == "" {
		b.sendMessage(chatID, "Укажите код валюты для получения курса\n\nПример: <code>/exchange USD</code>\n\nДоступно: USD, EUR, CNY, GBP, JPY и другие")
		return
	}

	b.sendMessage(chatID, "Получаю актуальный курс валют...")

	rateInfo, err := api.GetExchangeRate(currency)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("<b>Ошибка:</b> %s", err.Error()))
		return
	}

	b.sendMessage(chatID, rateInfo)
}
