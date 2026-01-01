package main

import (
	"mini-app-backend/internal/bot"
	"mini-app-backend/internal/config"
	"mini-app-backend/internal/logger"
	"mini-app-backend/internal/server"
	"net/http"
	"sync"
)

func main() {
	// Initialize logger
	logger.SetLogger(logger.NewDefaultLogger())

	cfg := config.Load()

	if cfg.TelegramBotToken == "" {
		logger.GetLogger().Fatal("❌ TELEGRAM_BOT_TOKEN не установлен")
	}

	tgBot, err := bot.New(cfg)
	if err != nil {
		logger.GetLogger().Fatalf("Failed created bot: %v", err)
	}

	httpServer := server.New(cfg)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		logger.GetLogger().Info("🤖 Start Telegram bot...")
		tgBot.Start()
	}()

	go func() {
		defer wg.Done()
		logger.GetLogger().Info("🌐 Start HTTP server...")
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Fatalf("Failed start HTTP server: %v", err)
		}
	}()

	logger.GetLogger().Info("✅ All services started")
	wg.Wait()
}
