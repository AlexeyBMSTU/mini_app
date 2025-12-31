# Этап сборки
FROM golang:1.25-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем Air для горячей перезагрузки
RUN go install github.com/air-verse/air@latest

# Копируем go mod файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Создаем директорию для временных файлов Air
RUN mkdir -p ./tmp

# Финальный этап для сервера
FROM golang:1.25-alpine AS server

# Устанавливаем ca-certificates для HTTPS и обновляем пакеты
RUN apk --no-cache add ca-certificates tzdata git

# Копируем Air из этапа сборки
COPY --from=builder /go/bin/air /go/bin/air

# Создаем директорию для приложения
WORKDIR /app

# Копируем конфигурацию Air
COPY --from=builder /app/.air.toml ./

# Экспонируем порт
EXPOSE 8080

# Запускаем Air для сервера
CMD ["air", "-c", ".air.toml"]

# Финальный этап для бота
FROM golang:1.25-alpine AS bot

# Устанавливаем ca-certificates для HTTPS и обновляем пакеты
RUN apk --no-cache add ca-certificates tzdata git

# Копируем Air из этапа сборки
COPY --from=builder /go/bin/air /go/bin/air

# Создаем директорию для приложения
WORKDIR /app

# Копируем конфигурацию Air для бота
COPY --from=builder /app/.air.bot.toml ./.air.toml

# Запускаем Air для бота
CMD ["air", "-c", ".air.toml"]