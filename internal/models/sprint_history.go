package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"sprint-backlog/pkg/constants"
)

type SprintHistory struct {
	ID        uuid.UUID              `gorm:"type:uuid;primary_key" json:"id"`
	SprintID  uuid.UUID              `gorm:"type:uuid;not null;index" json:"sprint_id"`
	UserID    uuid.UUID              `gorm:"type:uuid;not null" json:"user_id"`
	ItemID    *uuid.UUID             `gorm:"type:uuid" json:"item_id"`
	Action    constants.SprintAction `gorm:"type:varchar(50);not null" json:"action"`
	OldValue  datatypes.JSON         `gorm:"type:jsonb" json:"old_value"`
	NewValue  datatypes.JSON         `gorm:"type:jsonb" json:"new_value"`
	Timestamp time.Time              `gorm:"not null;default:now()" json:"timestamp"`

	// Relations
	Sprint Sprint       `gorm:"foreignKey:SprintID" json:"sprint,omitempty"`
	User   User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Item   *BacklogItem `gorm:"foreignKey:ItemID" json:"item,omitempty"`
}

func (h *SprintHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	if h.Timestamp.IsZero() {
		h.Timestamp = time.Now()
	}
	return nil
}

// TableName specifies the table name for SprintHistory model
func (SprintHistory) TableName() string {
	return "sprint_histories"
}
