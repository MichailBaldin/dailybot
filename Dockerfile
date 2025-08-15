FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o dailybot cmd/bot/main.go

# Проверяем, что файл создался
RUN ls -la /app/dailybot

# Финальный образ
FROM alpine:latest

# Устанавливаем только необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем исполняемый файл
COPY --from=builder /app/dailybot ./dailybot

# Проверяем, что файл скопировался
RUN ls -la /app/dailybot

# Даем права на выполнение
RUN chmod +x /app/dailybot

# Открываем порт для админки
EXPOSE 8080

# Запускаем приложение
CMD ["./dailybot"]