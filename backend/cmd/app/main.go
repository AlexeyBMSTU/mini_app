package main

import (
	"log"
	"net/http"
	"sync"

	"mini-app-backend/internal/bot"
	"mini-app-backend/internal/config"
	"mini-app-backend/internal/server"
)

func main() {
	cfg := config.Load()

	if cfg.TelegramBotToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN не установлен")
	}

	tgBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed created bot: %v", err)
	}

	httpServer := server.New(cfg)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Println("🤖 Start Telegram bot...")
		tgBot.Start()
	}()

	go func() {
		defer wg.Done()
		log.Println("🌐 Start HTTP server...")
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed start HTTP server: %v", err)
		}
	}()

	log.Println("✅ All services started")
	wg.Wait()
}
