package bot

import (
	"dailybot/internal/api"
	"fmt"
	"strings"
)

func (b *Bot) handleStart(chatID int64) {
	text := `<b>Привет! Я ДейлиБот - твой помощник на каждый день!</b>

<b>Мои команды:</b>
/help - подробная справка
/weather [город] - узнать погоду
/news - узнать актуальные новости`

	b.sendMessage(chatID, text)
}

func (b *Bot) handleHelp(chatID int64) {
	text := `<b>Справка по командам:</b>

<b>Мои команды:</b>
<b>/weather [город]</b> - получить прогноз погоды
Пример: <code>/weather Москва</code>

<b>/news</b> - топ-5 главных новостей дня
Актуальные новости из российских источников

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

	newsInfo, err := api.GetNews(b.config.NewsAPIKey)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("<b>Ошибка:</b> %s", err.Error()))
		return
	}

	b.sendMessage(chatID, newsInfo)
}
