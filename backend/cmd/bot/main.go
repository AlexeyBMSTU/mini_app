package main

import (
	"log"

	"mini-app-backend/internal/bot"
	"mini-app-backend/internal/config"
)

func main() {
	cfg := config.Load()

	if cfg.TelegramBotToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN not downloaded")
	}

	tgBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed created bot: %v", err)
	}

	log.Println("🤖 Start Telegram bot...")
	tgBot.Start()
}
