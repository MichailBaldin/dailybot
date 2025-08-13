# DailyBot - Telegram Bot на Go

Многофункциональный Telegram-бот для ежедневных задач, разработанный на языке Go с чистой архитектурой.

## Функциональность

- **🌤 Погода** - актуальная погода в любом городе (OpenWeather API)
- **💱 Курсы валют** - курсы валют по данным ЦБ РФ
- **📰 Новости** - главные новости дня (NewsAPI)

## Технологии

- Go 1.24+
- Telegram Bot API
- Clean Architecture
- Docker & Docker Compose
- HTML форматирование сообщений

## Быстрый запуск

```bash
# Клонируем репозиторий
git clone https://github.com/MichailBaldin/dailybot
cd dailybot

# Настраиваем переменные окружения
cp .env.example .env
# Отредактируйте .env файл с вашими токенами

# Запуск через Docker
docker-compose up --build

# Или локально
go run cmd/bot/main.go
```

## Конфигурация

Создайте .env файл:

```bash
TELEGRAM_BOT_TOKEN=your_bot_token_here
OPENWEATHER_API_KEY=your_openweather_key_here
NEWS_API_KEY=your_news_api_key_here
```

## Получение API ключей

Telegram Bot Token: @BotFather в Telegram
OpenWeather API: https://openweathermap.org/api (бесплатно)
News API: https://newsapi.org (опционально)

## Команды бота

/start - приветствие и список команд
/help - подробная справка
/weather [город] - прогноз погоды
/exchange [валюта] - курс валют
/news - главные новости

## Архитектура

```
cmd/bot/main.go          # Точка входа
internal/
├── config/              # Конфигурация
├── bot/                 # Логика бота
└── api/                 # Внешние API
    ├── weather.go       # OpenWeather API
    ├── exchange.go      # ЦБ РФ API
    └── news.go          # News API
```

## Docker

```bash
# Сборка
docker build -t dailybot .

# Запуск
docker run --env-file .env dailybot
```

## Лицензия

MIT License - используйте свободно в коммерческих и личных проектах.
