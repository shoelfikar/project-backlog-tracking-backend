package request

import (
	"github.com/google/uuid"

	"sprint-backlog/pkg/constants"
)

// CreateBacklogItemRequest represents the request body for creating a backlog item
type CreateBacklogItemRequest struct {
	ProjectID   uuid.UUID           `json:"project_id" binding:"required"`
	Title       string              `json:"title" binding:"required,min=1,max=200"`
	Description string              `json:"description" binding:"max=5000"`
	Type        constants.ItemType  `json:"type" binding:"required"`
	Priority    constants.Priority  `json:"priority" binding:"required"`
	Status      constants.ItemStatus `json:"status"`
	StoryPoints *int                `json:"story_points" binding:"omitempty,min=0,max=100"`
	Labels      []string            `json:"labels"`
	SprintID    *uuid.UUID          `json:"sprint_id"`
}

// UpdateBacklogItemRequest represents the request body for updating a backlog item
type UpdateBacklogItemRequest struct {
	Title       string              `json:"title" binding:"omitempty,min=1,max=200"`
	Description string              `json:"description" binding:"max=5000"`
	Type        constants.ItemType  `json:"type"`
	Priority    constants.Priority  `json:"priority"`
	Status      constants.ItemStatus `json:"status"`
	StoryPoints *int                `json:"story_points" binding:"omitempty,min=0,max=100"`
	Labels      []string            `json:"labels"`
	SprintID    *uuid.UUID          `json:"sprint_id"`
}

// UpdateStatusRequest represents the request body for updating item status
type UpdateStatusRequest struct {
	Status constants.ItemStatus `json:"status" binding:"required"`
}

// UpdatePriorityRequest represents the request body for updating item priority
type UpdatePriorityRequest struct {
	Priority constants.Priority `json:"priority" binding:"required"`
}

// AddLabelRequest represents the request body for adding a label
type AddLabelRequest struct {
	Label string `json:"label" binding:"required,min=1,max=50"`
}

// AddCommentRequest represents the request body for adding a comment
type AddCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=2000"`
}

// BacklogQueryParams represents query parameters for listing backlog items
type BacklogQueryParams struct {
	Search   string   `form:"search"`
	Type     []string `form:"type"`
	Priority []string `form:"priority"`
	Status   []string `form:"status"`
	SprintID string   `form:"sprint_id"`
	Labels   []string `form:"labels"`
	Page     int      `form:"page" binding:"omitempty,min=1"`
	Limit    int      `form:"limit" binding:"omitempty,min=1,max=100"`
}
