package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sprint-backlog/internal/dto/request"
	"sprint-backlog/internal/service"
	"sprint-backlog/internal/utils"
	"sprint-backlog/pkg/constants"
)

type BacklogHandler struct {
	backlogService service.BacklogService
}

func NewBacklogHandler(backlogService service.BacklogService) *BacklogHandler {
	return &BacklogHandler{
		backlogService: backlogService,
	}
}

// Create handles POST /api/backlog
// @Summary Create a new backlog item
// @Description Create a new backlog item
// @Tags backlog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateBacklogItemRequest true "Create backlog item request"
// @Success 201 {object} response.BacklogItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog [post]
func (h *BacklogHandler) Create(c *gin.Context) {
	var req request.CreateBacklogItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	item, err := h.backlogService.Create(&req, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidItemType):
			utils.RespondBadRequest(c, "Invalid item type", err.Error())
		case errors.Is(err, service.ErrInvalidPriority):
			utils.RespondBadRequest(c, "Invalid priority", err.Error())
		case errors.Is(err, service.ErrInvalidStatus):
			utils.RespondBadRequest(c, "Invalid status", err.Error())
		default:
			utils.RespondInternalError(c, "Failed to create backlog item", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, "Backlog item created successfully", item)
}

// GetAll handles GET /api/backlog
// @Summary Get all backlog items
// @Description Get all backlog items with optional filters and pagination
// @Tags backlog
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search in title and description"
// @Param type query []string false "Filter by type (Story, Task, Bug, Epic)"
// @Param priority query []string false "Filter by priority (Critical, High, Medium, Low)"
// @Param status query []string false "Filter by status (New, Ready, In Progress, In Review, Done)"
// @Param sprint_id query string false "Filter by sprint ID or 'none' for unassigned"
// @Param labels query []string false "Filter by labels"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BacklogListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog [get]
func (h *BacklogHandler) GetAll(c *gin.Context) {
	var params request.BacklogQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.RespondBadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.backlogService.GetAll(&params)
	if err != nil {
		utils.RespondInternalError(c, "Failed to fetch backlog items", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", result)
}

// GetByID handles GET /api/backlog/:id
// @Summary Get backlog item by ID
// @Description Get a backlog item by its ID
// @Tags backlog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Success 200 {object} response.BacklogItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /backlog/{id} [get]
func (h *BacklogHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	item, err := h.backlogService.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrBacklogItemNotFound) {
			utils.RespondNotFound(c, "Backlog item not found")
			return
		}
		utils.RespondInternalError(c, "Failed to fetch backlog item", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", item)
}

// Update handles PUT /api/backlog/:id
// @Summary Update a backlog item
// @Description Update an existing backlog item
// @Tags backlog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Param request body request.UpdateBacklogItemRequest true "Update backlog item request"
// @Success 200 {object} response.BacklogItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog/{id} [put]
func (h *BacklogHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	var req request.UpdateBacklogItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	item, err := h.backlogService.Update(id, &req, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBacklogItemNotFound):
			utils.RespondNotFound(c, "Backlog item not found")
		case errors.Is(err, service.ErrInvalidItemType):
			utils.RespondBadRequest(c, "Invalid item type", err.Error())
		case errors.Is(err, service.ErrInvalidPriority):
			utils.RespondBadRequest(c, "Invalid priority", err.Error())
		case errors.Is(err, service.ErrInvalidStatus):
			utils.RespondBadRequest(c, "Invalid status", err.Error())
		default:
			utils.RespondInternalError(c, "Failed to update backlog item", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Backlog item updated successfully", item)
}

// Delete handles DELETE /api/backlog/:id
// @Summary Delete a backlog item
// @Description Delete a backlog item by its ID
// @Tags backlog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog/{id} [delete]
func (h *BacklogHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	if err := h.backlogService.Delete(id); err != nil {
		if errors.Is(err, service.ErrBacklogItemNotFound) {
			utils.RespondNotFound(c, "Backlog item not found")
			return
		}
		utils.RespondInternalError(c, "Failed to delete backlog item", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Backlog item deleted successfully", nil)
}

// UpdateStatus handles PATCH /api/backlog/:id/status
// @Summary Update backlog item status
// @Description Update the status of a backlog item
// @Tags backlog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Param request body request.UpdateStatusRequest true "Update status request"
// @Success 200 {object} response.BacklogItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog/{id}/status [patch]
func (h *BacklogHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	var req request.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	item, err := h.backlogService.UpdateStatus(id, req.Status, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBacklogItemNotFound):
			utils.RespondNotFound(c, "Backlog item not found")
		case errors.Is(err, service.ErrInvalidStatus):
			utils.RespondBadRequest(c, "Invalid status", err.Error())
		default:
			utils.RespondInternalError(c, "Failed to update status", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Status updated successfully", item)
}

// UpdatePriority handles PATCH /api/backlog/:id/priority
// @Summary Update backlog item priority
// @Description Update the priority of a backlog item
// @Tags backlog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Param request body request.UpdatePriorityRequest true "Update priority request"
// @Success 200 {object} response.BacklogItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog/{id}/priority [patch]
func (h *BacklogHandler) UpdatePriority(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	var req request.UpdatePriorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	item, err := h.backlogService.UpdatePriority(id, req.Priority, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBacklogItemNotFound):
			utils.RespondNotFound(c, "Backlog item not found")
		case errors.Is(err, service.ErrInvalidPriority):
			utils.RespondBadRequest(c, "Invalid priority", err.Error())
		default:
			utils.RespondInternalError(c, "Failed to update priority", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Priority updated successfully", item)
}

// AddLabel handles POST /api/backlog/:id/labels
// @Summary Add a label to backlog item
// @Description Add a label to a backlog item
// @Tags backlog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Param request body request.AddLabelRequest true "Add label request"
// @Success 200 {object} response.BacklogItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog/{id}/labels [post]
func (h *BacklogHandler) AddLabel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	var req request.AddLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	item, err := h.backlogService.AddLabel(id, req.Label, userID)
	if err != nil {
		if errors.Is(err, service.ErrBacklogItemNotFound) {
			utils.RespondNotFound(c, "Backlog item not found")
			return
		}
		utils.RespondInternalError(c, "Failed to add label", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Label added successfully", item)
}

// RemoveLabel handles DELETE /api/backlog/:id/labels/:label
// @Summary Remove a label from backlog item
// @Description Remove a label from a backlog item
// @Tags backlog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Param label path string true "Label to remove"
// @Success 200 {object} response.BacklogItemResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog/{id}/labels/{label} [delete]
func (h *BacklogHandler) RemoveLabel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	label := c.Param("label")
	if label == "" {
		utils.RespondBadRequest(c, "Label is required", "")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	item, err := h.backlogService.RemoveLabel(id, label, userID)
	if err != nil {
		if errors.Is(err, service.ErrBacklogItemNotFound) {
			utils.RespondNotFound(c, "Backlog item not found")
			return
		}
		utils.RespondInternalError(c, "Failed to remove label", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Label removed successfully", item)
}

// AddComment handles POST /api/backlog/:id/comments
// @Summary Add a comment to backlog item
// @Description Add a comment to a backlog item
// @Tags backlog
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Param request body request.AddCommentRequest true "Add comment request"
// @Success 201 {object} response.ItemHistoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog/{id}/comments [post]
func (h *BacklogHandler) AddComment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	var req request.AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	history, err := h.backlogService.AddComment(id, req.Content, userID)
	if err != nil {
		if errors.Is(err, service.ErrBacklogItemNotFound) {
			utils.RespondNotFound(c, "Backlog item not found")
			return
		}
		utils.RespondInternalError(c, "Failed to add comment", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, "Comment added successfully", history)
}

// GetHistory handles GET /api/backlog/:id/history
// @Summary Get backlog item history
// @Description Get the history of a backlog item
// @Tags backlog
// @Produce json
// @Security BearerAuth
// @Param id path string true "Backlog Item ID"
// @Success 200 {array} response.ItemHistoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /backlog/{id}/history [get]
func (h *BacklogHandler) GetHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid backlog item ID", "ID must be a valid UUID")
		return
	}

	history, err := h.backlogService.GetHistory(id)
	if err != nil {
		if errors.Is(err, service.ErrBacklogItemNotFound) {
			utils.RespondNotFound(c, "Backlog item not found")
			return
		}
		utils.RespondInternalError(c, "Failed to fetch history", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", history)
}

// Helper to convert string to constants.ItemStatus
func parseStatus(s string) constants.ItemStatus {
	return constants.ItemStatus(s)
}

// Helper to convert string to constants.Priority
func parsePriority(s string) constants.Priority {
	return constants.Priority(s)
}
