package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sprint-backlog/internal/models"
	"sprint-backlog/pkg/constants"
)

type SprintRepository interface {
	Create(sprint *models.Sprint) error
	GetByID(id uuid.UUID) (*models.Sprint, error)
	GetByProjectID(projectID uuid.UUID, filters SprintFilters) ([]models.Sprint, int64, error)
	GetAll(filters SprintFilters) ([]models.Sprint, int64, error)
	GetActive(projectID uuid.UUID) (*models.Sprint, error)
	Update(sprint *models.Sprint) error
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status constants.SprintStatus) error
	GetItemsBySprintID(sprintID uuid.UUID) ([]models.BacklogItem, error)
	CalculateVelocity(sprintID uuid.UUID) (int, error)
}

type SprintFilters struct {
	ProjectID *uuid.UUID
	Status    []constants.SprintStatus
	Page      int
	Limit     int
}

type sprintRepository struct {
	db *gorm.DB
}

func NewSprintRepository(db *gorm.DB) SprintRepository {
	return &sprintRepository{db: db}
}

func (r *sprintRepository) Create(sprint *models.Sprint) error {
	return r.db.Create(sprint).Error
}

func (r *sprintRepository) GetByID(id uuid.UUID) (*models.Sprint, error) {
	var sprint models.Sprint
	err := r.db.Preload("CreatedBy").Preload("Project").
		Where("id = ?", id).First(&sprint).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sprint, nil
}

func (r *sprintRepository) GetByProjectID(projectID uuid.UUID, filters SprintFilters) ([]models.Sprint, int64, error) {
	var sprints []models.Sprint
	var total int64

	query := r.db.Model(&models.Sprint{}).Where("project_id = ?", projectID)
	query = r.applyFilters(query, filters)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get results
	if filters.Page > 0 && filters.Limit > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset).Limit(filters.Limit)
	}

	err := query.Preload("CreatedBy").
		Order("start_date DESC").
		Find(&sprints).Error

	return sprints, total, err
}

func (r *sprintRepository) GetAll(filters SprintFilters) ([]models.Sprint, int64, error) {
	var sprints []models.Sprint
	var total int64

	query := r.db.Model(&models.Sprint{})
	query = r.applyFilters(query, filters)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and get results
	if filters.Page > 0 && filters.Limit > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset).Limit(filters.Limit)
	}

	err := query.Preload("CreatedBy").Preload("Project").
		Order("start_date DESC").
		Find(&sprints).Error

	return sprints, total, err
}

func (r *sprintRepository) GetActive(projectID uuid.UUID) (*models.Sprint, error) {
	var sprint models.Sprint
	err := r.db.Preload("CreatedBy").Preload("Project").
		Where("project_id = ? AND status = ?", projectID, constants.SprintStatusActive).
		First(&sprint).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sprint, nil
}

func (r *sprintRepository) Update(sprint *models.Sprint) error {
	return r.db.Save(sprint).Error
}

func (r *sprintRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Sprint{}, "id = ?", id).Error
}

func (r *sprintRepository) UpdateStatus(id uuid.UUID, status constants.SprintStatus) error {
	return r.db.Model(&models.Sprint{}).Where("id = ?", id).Update("status", status).Error
}

func (r *sprintRepository) GetItemsBySprintID(sprintID uuid.UUID) ([]models.BacklogItem, error) {
	var items []models.BacklogItem
	err := r.db.Preload("CreatedBy").
		Where("sprint_id = ?", sprintID).
		Order("position ASC").
		Find(&items).Error
	return items, err
}

func (r *sprintRepository) CalculateVelocity(sprintID uuid.UUID) (int, error) {
	var velocity int
	err := r.db.Model(&models.BacklogItem{}).
		Where("sprint_id = ? AND status = ?", sprintID, constants.ItemStatusDone).
		Select("COALESCE(SUM(story_points), 0)").
		Scan(&velocity).Error
	return velocity, err
}

func (r *sprintRepository) applyFilters(query *gorm.DB, filters SprintFilters) *gorm.DB {
	// Project filter
	if filters.ProjectID != nil {
		query = query.Where("project_id = ?", *filters.ProjectID)
	}

	// Status filter
	if len(filters.Status) > 0 {
		query = query.Where("status IN ?", filters.Status)
	}

	return query
}
