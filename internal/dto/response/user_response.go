package response

import (
	"time"

	"github.com/google/uuid"
)

// UserActivityResponse represents a unified activity entry (from both item and sprint history)
type UserActivityResponse struct {
	ID           uuid.UUID           `json:"id"`
	Type         string              `json:"type"` // "item" or "sprint"
	Action       string              `json:"action"`
	Timestamp    time.Time           `json:"timestamp"`
	FieldChanged *string             `json:"field_changed,omitempty"`
	OldValue     interface{}         `json:"old_value,omitempty"`
	NewValue     interface{}         `json:"new_value,omitempty"`
	Comment      *string             `json:"comment,omitempty"`
	Item         *BacklogItemSummary `json:"item,omitempty"`
	Sprint       *SprintSummaryFull  `json:"sprint,omitempty"`
}

// SprintSummaryFull represents a sprint summary with more details
type SprintSummaryFull struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// UserActivitiesResponse represents a paginated list of user activities
type UserActivitiesResponse struct {
	Activities []UserActivityResponse `json:"activities"`
	Total      int                    `json:"total"`
	Limit      int                    `json:"limit"`
}
