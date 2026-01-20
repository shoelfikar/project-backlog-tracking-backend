package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"sprint-backlog/internal/models"
	"sprint-backlog/pkg/constants"
)

type BacklogRepository interface {
	Create(item *models.BacklogItem) error
	GetByID(id uuid.UUID) (*models.BacklogItem, error)
	GetByProjectID(projectID uuid.UUID, filters BacklogFilters) ([]models.BacklogItem, int64, error)
	GetBySprintID(sprintID uuid.UUID) ([]models.BacklogItem, error)
	GetAll(filters BacklogFilters) ([]models.BacklogItem, int64, error)
	Update(item *models.BacklogItem) error
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, status constants.ItemStatus) error
	UpdatePriority(id uuid.UUID, priority constants.Priority) error
	AddLabel(id uuid.UUID, label string) error
	RemoveLabel(id uuid.UUID, label string) error
	GetMaxPosition(projectID uuid.UUID) (int, error)
}

type BacklogFilters struct {
	Search    string
	Type      []constants.ItemType
	Priority  []constants.Priority
	Status    []constants.ItemStatus
	SprintID  *uuid.UUID
	Labels    []string
	Page      int
	Limit     int
}

type backlogRepository struct {
	db *gorm.DB
}

func NewBacklogRepository(db *gorm.DB) BacklogRepository {
	return &backlogRepository{db: db}
}

func (r *backlogRepository) Create(item *models.BacklogItem) error {
	return r.db.Create(item).Error
}

func (r *backlogRepository) GetByID(id uuid.UUID) (*models.BacklogItem, error) {
	var item models.BacklogItem
	err := r.db.Preload("CreatedBy").Preload("Sprint").Preload("Project").
		Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *backlogRepository) GetByProjectID(projectID uuid.UUID, filters BacklogFilters) ([]models.BacklogItem, int64, error) {
	var items []models.BacklogItem
	var total int64

	query := r.db.Model(&models.BacklogItem{}).Where("project_id = ?", projectID)
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

	err := query.Preload("CreatedBy").Preload("Sprint").
		Order("position ASC, created_at DESC").
		Find(&items).Error

	return items, total, err
}

func (r *backlogRepository) GetBySprintID(sprintID uuid.UUID) ([]models.BacklogItem, error) {
	var items []models.BacklogItem
	err := r.db.Preload("CreatedBy").
		Where("sprint_id = ?", sprintID).
		Order("position ASC").
		Find(&items).Error
	return items, err
}

func (r *backlogRepository) GetAll(filters BacklogFilters) ([]models.BacklogItem, int64, error) {
	var items []models.BacklogItem
	var total int64

	query := r.db.Model(&models.BacklogItem{})
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

	err := query.Preload("CreatedBy").Preload("Sprint").Preload("Project").
		Order("created_at DESC").
		Find(&items).Error

	return items, total, err
}

func (r *backlogRepository) Update(item *models.BacklogItem) error {
	return r.db.Save(item).Error
}

func (r *backlogRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BacklogItem{}, "id = ?", id).Error
}

func (r *backlogRepository) UpdateStatus(id uuid.UUID, status constants.ItemStatus) error {
	return r.db.Model(&models.BacklogItem{}).Where("id = ?", id).Update("status", status).Error
}

func (r *backlogRepository) UpdatePriority(id uuid.UUID, priority constants.Priority) error {
	return r.db.Model(&models.BacklogItem{}).Where("id = ?", id).Update("priority", priority).Error
}

func (r *backlogRepository) AddLabel(id uuid.UUID, label string) error {
	return r.db.Exec(
		"UPDATE backlog_items SET labels = array_append(labels, ?) WHERE id = ? AND NOT (? = ANY(labels))",
		label, id, label,
	).Error
}

func (r *backlogRepository) RemoveLabel(id uuid.UUID, label string) error {
	return r.db.Exec(
		"UPDATE backlog_items SET labels = array_remove(labels, ?) WHERE id = ?",
		label, id,
	).Error
}

func (r *backlogRepository) GetMaxPosition(projectID uuid.UUID) (int, error) {
	var maxPosition int
	err := r.db.Model(&models.BacklogItem{}).
		Where("project_id = ?", projectID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *backlogRepository) applyFilters(query *gorm.DB, filters BacklogFilters) *gorm.DB {
	// Search filter
	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	// Type filter
	if len(filters.Type) > 0 {
		query = query.Where("type IN ?", filters.Type)
	}

	// Priority filter
	if len(filters.Priority) > 0 {
		query = query.Where("priority IN ?", filters.Priority)
	}

	// Status filter
	if len(filters.Status) > 0 {
		query = query.Where("status IN ?", filters.Status)
	}

	// Sprint filter
	if filters.SprintID != nil {
		if *filters.SprintID == uuid.Nil {
			query = query.Where("sprint_id IS NULL")
		} else {
			query = query.Where("sprint_id = ?", *filters.SprintID)
		}
	}

	// Labels filter
	if len(filters.Labels) > 0 {
		query = query.Where("labels && ?", pq.Array(filters.Labels))
	}

	return query
}
