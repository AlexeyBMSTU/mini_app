package main

import (
	"log"
	"net/http"

	"mini-app-backend/internal/config"
	"mini-app-backend/internal/server"
)

func main() {
	cfg := config.Load()

	httpServer := server.New(cfg)

	log.Println("🌐 Start HTTP server...")
	if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed start HTTP server: %v", err)
	}
}
