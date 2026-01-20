package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sprint-backlog/internal/models"
)

type ItemHistoryRepository interface {
	Create(history *models.ItemHistory) error
	GetByItemID(itemID uuid.UUID) ([]models.ItemHistory, error)
	GetByUserID(userID uuid.UUID, limit int) ([]models.ItemHistory, error)
}

type itemHistoryRepository struct {
	db *gorm.DB
}

func NewItemHistoryRepository(db *gorm.DB) ItemHistoryRepository {
	return &itemHistoryRepository{db: db}
}

func (r *itemHistoryRepository) Create(history *models.ItemHistory) error {
	return r.db.Create(history).Error
}

func (r *itemHistoryRepository) GetByItemID(itemID uuid.UUID) ([]models.ItemHistory, error) {
	var histories []models.ItemHistory
	err := r.db.Preload("User").
		Where("item_id = ?", itemID).
		Order("timestamp DESC").
		Find(&histories).Error
	return histories, err
}

func (r *itemHistoryRepository) GetByUserID(userID uuid.UUID, limit int) ([]models.ItemHistory, error) {
	var histories []models.ItemHistory
	query := r.db.Preload("User").Preload("Item").
		Where("user_id = ?", userID).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&histories).Error
	return histories, err
}
