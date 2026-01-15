package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Key         string         `gorm:"uniqueIndex;not null;size:10" json:"key"`
	Description *string        `json:"description"`
	CreatedByID uuid.UUID      `gorm:"type:uuid;not null" json:"created_by_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	CreatedBy    User          `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	Sprints      []Sprint      `gorm:"foreignKey:ProjectID" json:"sprints,omitempty"`
	BacklogItems []BacklogItem `gorm:"foreignKey:ProjectID" json:"backlog_items,omitempty"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Project model
func (Project) TableName() string {
	return "projects"
}
