package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"sprint-backlog/pkg/constants"
)

type Sprint struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key" json:"id"`
	ProjectID   uuid.UUID              `gorm:"type:uuid;not null;index" json:"project_id"`
	CreatedByID uuid.UUID              `gorm:"type:uuid;not null" json:"created_by_id"`
	Name        string                 `gorm:"not null" json:"name"`
	Goal        *string                `json:"goal"`
	StartDate   time.Time              `gorm:"not null" json:"start_date"`
	EndDate     time.Time              `gorm:"not null" json:"end_date"`
	Status      constants.SprintStatus `gorm:"type:varchar(20);not null;default:'Planning'" json:"status"`
	Velocity    *int                   `json:"velocity"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `gorm:"index" json:"-"`

	// Relations
	Project   Project         `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	CreatedBy User            `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	Items     []BacklogItem   `gorm:"foreignKey:SprintID" json:"items,omitempty"`
	History   []SprintHistory `gorm:"foreignKey:SprintID" json:"history,omitempty"`
}

func (s *Sprint) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for Sprint model
func (Sprint) TableName() string {
	return "sprints"
}
