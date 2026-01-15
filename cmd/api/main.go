package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"sprint-backlog/internal/config"
	"sprint-backlog/internal/database"
	"sprint-backlog/internal/router"
)

func main() {
	// Load configuration
	config.Load()
	log.Println("Configuration loaded")

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Connect to database
	database.Connect()

	// Run migrations
	database.RunMigrations()

	// Setup router
	r := router.Setup(database.DB)

	// Start server
	addr := ":" + config.AppConfig.Port
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
