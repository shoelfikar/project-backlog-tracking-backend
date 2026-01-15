package database

import (
	"log"

	"sprint-backlog/internal/models"
)

func RunMigrations() {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.BacklogItem{},
		&models.Sprint{},
		&models.ItemHistory{},
		&models.SprintHistory{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")
}
