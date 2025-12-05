package main

import (
	"log"

	"mini-app-backend/internal/bot"
	"mini-app-backend/internal/config"
)

func main() {
	cfg := config.Load()

	if cfg.TelegramBotToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN не установлен")
	}

	tgBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	tgBot.Start()
}