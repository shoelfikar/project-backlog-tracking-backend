package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sprint-backlog/internal/models"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	GetByID(id uuid.UUID) (*models.Project, error)
	GetByKey(key string) (*models.Project, error)
	GetAll() ([]models.Project, error)
	GetAllWithPagination(page, limit int) ([]models.Project, int64, error)
	Update(project *models.Project) error
	Delete(id uuid.UUID) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) GetByID(id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("CreatedBy").Where("id = ?", id).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetByKey(key string) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("CreatedBy").Where("key = ?", key).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *projectRepository) GetAll() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Preload("CreatedBy").Order("created_at DESC").Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetAllWithPagination(page, limit int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	// Count total
	if err := r.db.Model(&models.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := r.db.Preload("CreatedBy").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&projects).Error

	return projects, total, err
}

func (r *projectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}
