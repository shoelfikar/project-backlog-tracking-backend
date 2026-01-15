package response

import (
	"time"

	"github.com/google/uuid"

	"sprint-backlog/internal/models"
)

// ProjectResponse represents a project in API responses
type ProjectResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Key         string        `json:"key"`
	Description string        `json:"description"`
	CreatedBy   *UserResponse `json:"created_by,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ProjectListResponse represents a paginated list of projects
type ProjectListResponse struct {
	Projects   []ProjectResponse `json:"projects"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// ToProjectResponse converts a Project model to ProjectResponse
func ToProjectResponse(project *models.Project) *ProjectResponse {
	if project == nil {
		return nil
	}

	resp := &ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		Key:       project.Key,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}

	// Handle nullable description
	if project.Description != nil {
		resp.Description = *project.Description
	}

	// Include CreatedBy if preloaded (check if ID is not zero)
	if project.CreatedBy.ID != uuid.Nil {
		resp.CreatedBy = ToUserResponse(&project.CreatedBy)
	}

	return resp
}

// ToProjectListResponse converts a slice of Project models to ProjectListResponse
func ToProjectListResponse(projects []models.Project, total int64, page, limit int) *ProjectListResponse {
	projectResponses := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = *ToProjectResponse(&project)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &ProjectListResponse{
		Projects:   projectResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
