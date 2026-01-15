package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sprint-backlog/internal/dto/request"
	"sprint-backlog/internal/service"
	"sprint-backlog/internal/utils"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// GetAll handles GET /api/projects
// @Summary Get all projects
// @Description Get all projects with optional pagination
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.ProjectListResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /projects [get]
func (h *ProjectHandler) GetAll(c *gin.Context) {
	var params request.ProjectQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.RespondBadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	// Set defaults
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	result, err := h.projectService.GetAllWithPagination(params.Page, params.Limit)
	if err != nil {
		utils.RespondInternalError(c, "Failed to fetch projects", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", result)
}

// Create handles POST /api/projects
// @Summary Create a new project
// @Description Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateProjectRequest true "Create project request"
// @Success 201 {object} response.ProjectResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var req request.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Get user ID from context
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	project, err := h.projectService.Create(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrProjectKeyExists) {
			utils.RespondError(c, http.StatusConflict, "Project key already exists", "PROJECT_KEY_EXISTS", "")
			return
		}
		utils.RespondInternalError(c, "Failed to create project", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, "Project created successfully", project)
}

// GetByID handles GET /api/projects/:id
// @Summary Get project by ID
// @Description Get a project by its ID
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} response.ProjectResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid project ID", "ID must be a valid UUID")
		return
	}

	project, err := h.projectService.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			utils.RespondNotFound(c, "Project not found")
			return
		}
		utils.RespondInternalError(c, "Failed to fetch project", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", project)
}

// Update handles PUT /api/projects/:id
// @Summary Update a project
// @Description Update an existing project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body request.UpdateProjectRequest true "Update project request"
// @Success 200 {object} response.ProjectResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid project ID", "ID must be a valid UUID")
		return
	}

	var req request.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	project, err := h.projectService.Update(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			utils.RespondNotFound(c, "Project not found")
			return
		}
		utils.RespondInternalError(c, "Failed to update project", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Project updated successfully", project)
}

// Delete handles DELETE /api/projects/:id
// @Summary Delete a project
// @Description Delete a project by its ID
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid project ID", "ID must be a valid UUID")
		return
	}

	if err := h.projectService.Delete(id); err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			utils.RespondNotFound(c, "Project not found")
			return
		}
		utils.RespondInternalError(c, "Failed to delete project", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Project deleted successfully", nil)
}
