package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sprint-backlog/internal/models"
)

type SprintHistoryRepository interface {
	Create(history *models.SprintHistory) error
	GetBySprintID(sprintID uuid.UUID) ([]models.SprintHistory, error)
	GetByUserID(userID uuid.UUID, limit int) ([]models.SprintHistory, error)
	GetAll(limit int) ([]models.SprintHistory, error)
}

type sprintHistoryRepository struct {
	db *gorm.DB
}

func NewSprintHistoryRepository(db *gorm.DB) SprintHistoryRepository {
	return &sprintHistoryRepository{db: db}
}

func (r *sprintHistoryRepository) Create(history *models.SprintHistory) error {
	return r.db.Create(history).Error
}

func (r *sprintHistoryRepository) GetBySprintID(sprintID uuid.UUID) ([]models.SprintHistory, error) {
	var histories []models.SprintHistory
	err := r.db.Preload("User").Preload("Item").
		Where("sprint_id = ?", sprintID).
		Order("timestamp DESC").
		Find(&histories).Error
	return histories, err
}

func (r *sprintHistoryRepository) GetByUserID(userID uuid.UUID, limit int) ([]models.SprintHistory, error) {
	var histories []models.SprintHistory
	query := r.db.Preload("User").Preload("Sprint").Preload("Item").
		Where("user_id = ?", userID).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&histories).Error
	return histories, err
}

func (r *sprintHistoryRepository) GetAll(limit int) ([]models.SprintHistory, error) {
	var histories []models.SprintHistory
	query := r.db.Preload("User").Preload("Sprint").Preload("Item").
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&histories).Error
	return histories, err
}
