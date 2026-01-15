package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"sprint-backlog/pkg/constants"
)

type ItemHistory struct {
	ID           uuid.UUID            `gorm:"type:uuid;primary_key" json:"id"`
	ItemID       uuid.UUID            `gorm:"type:uuid;not null;index" json:"item_id"`
	UserID       uuid.UUID            `gorm:"type:uuid;not null" json:"user_id"`
	Action       constants.ItemAction `gorm:"type:varchar(50);not null" json:"action"`
	FieldChanged *string              `json:"field_changed"`
	OldValue     datatypes.JSON       `gorm:"type:jsonb" json:"old_value"`
	NewValue     datatypes.JSON       `gorm:"type:jsonb" json:"new_value"`
	Comment      *string              `json:"comment"`
	Timestamp    time.Time            `gorm:"not null;default:now()" json:"timestamp"`

	// Relations
	Item BacklogItem `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	User User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (h *ItemHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	if h.Timestamp.IsZero() {
		h.Timestamp = time.Now()
	}
	return nil
}

// TableName specifies the table name for ItemHistory model
func (ItemHistory) TableName() string {
	return "item_histories"
}
