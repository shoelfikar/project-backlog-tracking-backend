package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"sprint-backlog/pkg/constants"
)

type BacklogItem struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID   uuid.UUID              `gorm:"type:uuid;not null;index" json:"project_id"`
	SprintID    *uuid.UUID             `gorm:"type:uuid;index" json:"sprint_id"`
	CreatedByID uuid.UUID              `gorm:"type:uuid;not null" json:"created_by_id"`
	Title       string                 `gorm:"not null" json:"title"`
	Description *string                `json:"description"`
	Type        constants.ItemType     `gorm:"type:varchar(20);not null;default:'Task'" json:"type"`
	Priority    constants.Priority     `gorm:"type:varchar(20);not null;default:'Medium'" json:"priority"`
	Status      constants.ItemStatus   `gorm:"type:varchar(20);not null;default:'New'" json:"status"`
	StoryPoints *int                   `json:"story_points"`
	Labels      pq.StringArray         `gorm:"type:text[]" json:"labels"`
	Position    int                    `gorm:"not null;default:0" json:"position"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `gorm:"index" json:"-"`

	// Relations
	Project   Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Sprint    *Sprint       `gorm:"foreignKey:SprintID" json:"sprint,omitempty"`
	CreatedBy User          `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	History   []ItemHistory `gorm:"foreignKey:ItemID" json:"history,omitempty"`
}

func (b *BacklogItem) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for BacklogItem model
func (BacklogItem) TableName() string {
	return "backlog_items"
}
