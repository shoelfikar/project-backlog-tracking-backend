package response

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"sprint-backlog/internal/models"
	"sprint-backlog/pkg/constants"
)

// SprintResponse represents a sprint in API responses
type SprintResponse struct {
	ID        uuid.UUID              `json:"id"`
	ProjectID uuid.UUID              `json:"project_id"`
	Name      string                 `json:"name"`
	Goal      string                 `json:"goal"`
	StartDate time.Time              `json:"start_date"`
	EndDate   time.Time              `json:"end_date"`
	Status    constants.SprintStatus `json:"status"`
	Velocity  *int                   `json:"velocity"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	CreatedBy *UserResponse          `json:"created_by,omitempty"`
	Project   *ProjectSummary        `json:"project,omitempty"`
}

// ProjectSummary represents a project summary in responses
type ProjectSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Key  string    `json:"key"`
}

// SprintWithItemsResponse represents a sprint with its items
type SprintWithItemsResponse struct {
	SprintResponse
	Items      []BacklogItemResponse `json:"items"`
	TotalItems int                   `json:"total_items"`
	TotalPoints int                  `json:"total_points"`
}

// SprintListResponse represents a paginated list of sprints
type SprintListResponse struct {
	Sprints    []SprintResponse `json:"sprints"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// SprintHistoryResponse represents a sprint history entry in API responses
type SprintHistoryResponse struct {
	ID        uuid.UUID              `json:"id"`
	SprintID  uuid.UUID              `json:"sprint_id"`
	UserID    uuid.UUID              `json:"user_id"`
	ItemID    *uuid.UUID             `json:"item_id"`
	Action    constants.SprintAction `json:"action"`
	OldValue  interface{}            `json:"old_value"`
	NewValue  interface{}            `json:"new_value"`
	Timestamp time.Time              `json:"timestamp"`
	User      *UserResponse          `json:"user,omitempty"`
	Item      *BacklogItemSummary    `json:"item,omitempty"`
}

// BacklogItemSummary represents a backlog item summary in responses
type BacklogItemSummary struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Type  string    `json:"type"`
}

// SprintReportResponse represents a sprint report
type SprintReportResponse struct {
	Sprint              SprintResponse `json:"sprint"`
	TotalItems          int            `json:"total_items"`
	CompletedItems      int            `json:"completed_items"`
	TotalStoryPoints    int            `json:"total_story_points"`
	CompletedStoryPoints int           `json:"completed_story_points"`
	Velocity            int            `json:"velocity"`
	CompletionPercentage float64       `json:"completion_percentage"`
}

// ToSprintResponse converts a Sprint model to SprintResponse
func ToSprintResponse(sprint *models.Sprint) *SprintResponse {
	if sprint == nil {
		return nil
	}

	resp := &SprintResponse{
		ID:        sprint.ID,
		ProjectID: sprint.ProjectID,
		Name:      sprint.Name,
		StartDate: sprint.StartDate,
		EndDate:   sprint.EndDate,
		Status:    sprint.Status,
		Velocity:  sprint.Velocity,
		CreatedAt: sprint.CreatedAt,
		UpdatedAt: sprint.UpdatedAt,
	}

	// Handle nullable goal
	if sprint.Goal != nil {
		resp.Goal = *sprint.Goal
	}

	// Include CreatedBy if preloaded
	if sprint.CreatedBy.ID != uuid.Nil {
		resp.CreatedBy = ToUserResponse(&sprint.CreatedBy)
	}

	// Include Project summary if preloaded
	if sprint.Project.ID != uuid.Nil {
		resp.Project = &ProjectSummary{
			ID:   sprint.Project.ID,
			Name: sprint.Project.Name,
			Key:  sprint.Project.Key,
		}
	}

	return resp
}

// ToSprintListResponse converts a slice of Sprint models to SprintListResponse
func ToSprintListResponse(sprints []models.Sprint, total int64, page, limit int) *SprintListResponse {
	sprintResponses := make([]SprintResponse, len(sprints))
	for i, sprint := range sprints {
		sprintResponses[i] = *ToSprintResponse(&sprint)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &SprintListResponse{
		Sprints:    sprintResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

// ToSprintWithItemsResponse converts a Sprint model with items
func ToSprintWithItemsResponse(sprint *models.Sprint, items []models.BacklogItem) *SprintWithItemsResponse {
	if sprint == nil {
		return nil
	}

	resp := &SprintWithItemsResponse{
		SprintResponse: *ToSprintResponse(sprint),
		Items:          make([]BacklogItemResponse, len(items)),
		TotalItems:     len(items),
		TotalPoints:    0,
	}

	for i, item := range items {
		resp.Items[i] = *ToBacklogItemResponse(&item)
		if item.StoryPoints != nil {
			resp.TotalPoints += *item.StoryPoints
		}
	}

	return resp
}

// ToSprintHistoryResponse converts a SprintHistory model to SprintHistoryResponse
func ToSprintHistoryResponse(history *models.SprintHistory) *SprintHistoryResponse {
	if history == nil {
		return nil
	}

	resp := &SprintHistoryResponse{
		ID:        history.ID,
		SprintID:  history.SprintID,
		UserID:    history.UserID,
		ItemID:    history.ItemID,
		Action:    history.Action,
		Timestamp: history.Timestamp,
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

	// Include Item summary if preloaded
	if history.Item != nil && history.Item.ID != uuid.Nil {
		resp.Item = &BacklogItemSummary{
			ID:    history.Item.ID,
			Title: history.Item.Title,
			Type:  string(history.Item.Type),
		}
	}

	return resp
}

// ToSprintHistoryListResponse converts a slice of SprintHistory models to responses
func ToSprintHistoryListResponse(histories []models.SprintHistory) []SprintHistoryResponse {
	responses := make([]SprintHistoryResponse, len(histories))
	for i, h := range histories {
		responses[i] = *ToSprintHistoryResponse(&h)
	}
	return responses
}
