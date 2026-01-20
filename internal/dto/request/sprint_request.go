package request

import (
	"time"

	"github.com/google/uuid"
)

// CreateSprintRequest represents the request body for creating a sprint
type CreateSprintRequest struct {
	ProjectID uuid.UUID  `json:"project_id" binding:"required"`
	Name      string     `json:"name" binding:"required,min=1,max=100"`
	Goal      string     `json:"goal" binding:"max=500"`
	StartDate time.Time  `json:"start_date" binding:"required"`
	EndDate   time.Time  `json:"end_date" binding:"required"`
}

// UpdateSprintRequest represents the request body for updating a sprint
type UpdateSprintRequest struct {
	Name      string    `json:"name" binding:"omitempty,min=1,max=100"`
	Goal      string    `json:"goal" binding:"max=500"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// AddItemToSprintRequest represents the request body for adding an item to a sprint
type AddItemToSprintRequest struct {
	ItemID uuid.UUID `json:"item_id" binding:"required"`
}

// SprintQueryParams represents query parameters for listing sprints
type SprintQueryParams struct {
	ProjectID string   `form:"project_id"`
	Status    []string `form:"status"`
	Page      int      `form:"page" binding:"omitempty,min=1"`
	Limit     int      `form:"limit" binding:"omitempty,min=1,max=100"`
}
