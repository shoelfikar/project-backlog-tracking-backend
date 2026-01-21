package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"sprint-backlog/internal/config"
	"sprint-backlog/internal/database"
	"sprint-backlog/internal/router"

	_ "sprint-backlog/docs" // Swagger docs
)

// @title Sprint Backlog API
// @version 1.0
// @description API for Sprint Backlog - A project management tool for agile teams
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@sprintbacklog.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header using Bearer scheme. Example: "Bearer {token}"

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
