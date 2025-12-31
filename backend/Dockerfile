# Этап сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл сервера
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main-server ./cmd/server/main.go

# Собираем бинарный файл бота
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main-bot ./cmd/bot/main.go

# Финальный этап для сервера
FROM alpine:latest AS server

# Устанавливаем ca-certificates для HTTPS и обновляем пакеты
RUN apk --no-cache add ca-certificates tzdata

# Создаем директорию для приложения
WORKDIR /root/

# Копируем бинарный файл сервера из этапа сборки
COPY --from=builder /app/main-server .

# Копируем файлы конфигурации
COPY --from=builder /app/.env .

# Экспонируем порт
EXPOSE 8080

# Запускаем сервер
CMD ["./main-server"]

# Финальный этап для бота
FROM alpine:latest AS bot

# Устанавливаем ca-certificates для HTTPS и обновляем пакеты
RUN apk --no-cache add ca-certificates tzdata

# Создаем директорию для приложения
WORKDIR /root/

# Копируем бинарный файл бота из этапа сборки
COPY --from=builder /app/main-bot .

# Копируем файлы конфигурации
COPY --from=builder /app/.env .

# Запускаем бота
CMD ["./main-bot"]