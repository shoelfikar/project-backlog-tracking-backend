package request

// CreateProjectRequest represents the request body for creating a project
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Key         string `json:"key" binding:"required,min=2,max=10,uppercase"`
	Description string `json:"description" binding:"max=500"`
}

// UpdateProjectRequest represents the request body for updating a project
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// ProjectQueryParams represents query parameters for listing projects
type ProjectQueryParams struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}
