package response

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"sprint-backlog/internal/models"
	"sprint-backlog/pkg/constants"
)

// BacklogItemResponse represents a backlog item in API responses
type BacklogItemResponse struct {
	ID          uuid.UUID            `json:"id"`
	ProjectID   uuid.UUID            `json:"project_id"`
	SprintID    *uuid.UUID           `json:"sprint_id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Type        constants.ItemType   `json:"type"`
	Priority    constants.Priority   `json:"priority"`
	Status      constants.ItemStatus `json:"status"`
	StoryPoints *int                 `json:"story_points"`
	Labels      []string             `json:"labels"`
	Position    int                  `json:"position"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	CreatedBy   *UserResponse        `json:"created_by,omitempty"`
	Sprint      *SprintSummary       `json:"sprint,omitempty"`
}

// SprintSummary represents a sprint summary in responses
type SprintSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// BacklogListResponse represents a paginated list of backlog items
type BacklogListResponse struct {
	Items      []BacklogItemResponse `json:"items"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

// ItemHistoryResponse represents an item history entry in API responses
type ItemHistoryResponse struct {
	ID           uuid.UUID            `json:"id"`
	ItemID       uuid.UUID            `json:"item_id"`
	UserID       uuid.UUID            `json:"user_id"`
	Action       constants.ItemAction `json:"action"`
	FieldChanged *string              `json:"field_changed"`
	OldValue     interface{}          `json:"old_value"`
	NewValue     interface{}          `json:"new_value"`
	Comment      *string              `json:"comment"`
	Timestamp    time.Time            `json:"timestamp"`
	User         *UserResponse        `json:"user,omitempty"`
}

// ToBacklogItemResponse converts a BacklogItem model to BacklogItemResponse
func ToBacklogItemResponse(item *models.BacklogItem) *BacklogItemResponse {
	if item == nil {
		return nil
	}

	resp := &BacklogItemResponse{
		ID:          item.ID,
		ProjectID:   item.ProjectID,
		SprintID:    item.SprintID,
		Title:       item.Title,
		Type:        item.Type,
		Priority:    item.Priority,
		Status:      item.Status,
		StoryPoints: item.StoryPoints,
		Labels:      item.Labels,
		Position:    item.Position,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	// Handle nullable description
	if item.Description != nil {
		resp.Description = *item.Description
	}

	// Handle labels
	if resp.Labels == nil {
		resp.Labels = []string{}
	}

	// Include CreatedBy if preloaded
	if item.CreatedBy.ID != uuid.Nil {
		resp.CreatedBy = ToUserResponse(&item.CreatedBy)
	}

	// Include Sprint summary if preloaded
	if item.Sprint != nil && item.Sprint.ID != uuid.Nil {
		resp.Sprint = &SprintSummary{
			ID:   item.Sprint.ID,
			Name: item.Sprint.Name,
		}
	}

	return resp
}

// ToBacklogListResponse converts a slice of BacklogItem models to BacklogListResponse
func ToBacklogListResponse(items []models.BacklogItem, total int64, page, limit int) *BacklogListResponse {
	itemResponses := make([]BacklogItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = *ToBacklogItemResponse(&item)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &BacklogListResponse{
		Items:      itemResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// ToItemHistoryResponse converts an ItemHistory model to ItemHistoryResponse
func ToItemHistoryResponse(history *models.ItemHistory) *ItemHistoryResponse {
	if history == nil {
		return nil
	}

	resp := &ItemHistoryResponse{
		ID:           history.ID,
		ItemID:       history.ItemID,
		UserID:       history.UserID,
		Action:       history.Action,
		FieldChanged: history.FieldChanged,
		Comment:      history.Comment,
		Timestamp:    history.Timestamp,
	}

	// Parse JSON values
	if history.OldValue != nil {
		var oldVal interface{}
		if err := json.Unmarshal(history.OldValue, &oldVal); err == nil {
			resp.OldValue = oldVal
		}
	}
	if history.NewValue != nil {
		var newVal interface{}
		if err := json.Unmarshal(history.NewValue, &newVal); err == nil {
			resp.NewValue = newVal
		}
	}

	// Include User if preloaded
	if history.User.ID != uuid.Nil {
		resp.User = ToUserResponse(&history.User)
	}

	return resp
}

// ToItemHistoryListResponse converts a slice of ItemHistory models to responses
func ToItemHistoryListResponse(histories []models.ItemHistory) []ItemHistoryResponse {
	responses := make([]ItemHistoryResponse, len(histories))
	for i, h := range histories {
		responses[i] = *ToItemHistoryResponse(&h)
	}
	return responses
}
