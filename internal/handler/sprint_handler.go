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

type SprintHandler struct {
	sprintService service.SprintService
}

func NewSprintHandler(sprintService service.SprintService) *SprintHandler {
	return &SprintHandler{
		sprintService: sprintService,
	}
}

// Create handles POST /api/sprints
// @Summary Create a new sprint
// @Description Create a new sprint
// @Tags sprints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateSprintRequest true "Create sprint request"
// @Success 201 {object} response.SprintResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints [post]
func (h *SprintHandler) Create(c *gin.Context) {
	var req request.CreateSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	sprint, err := h.sprintService.Create(&req, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDateRange) {
			utils.RespondBadRequest(c, "Invalid date range", err.Error())
			return
		}
		utils.RespondInternalError(c, "Failed to create sprint", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, "Sprint created successfully", sprint)
}

// GetAll handles GET /api/sprints
// @Summary Get all sprints
// @Description Get all sprints with optional filters and pagination
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param project_id query string false "Filter by project ID"
// @Param status query []string false "Filter by status (Planning, Active, Completed, Cancelled)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.SprintListResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints [get]
func (h *SprintHandler) GetAll(c *gin.Context) {
	var params request.SprintQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.RespondBadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	result, err := h.sprintService.GetAll(&params)
	if err != nil {
		utils.RespondInternalError(c, "Failed to fetch sprints", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", result)
}

// GetByID handles GET /api/sprints/:id
// @Summary Get sprint by ID
// @Description Get a sprint by its ID
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Success 200 {object} response.SprintResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /sprints/{id} [get]
func (h *SprintHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	sprint, err := h.sprintService.GetWithItems(id)
	if err != nil {
		if errors.Is(err, service.ErrSprintNotFound) {
			utils.RespondNotFound(c, "Sprint not found")
			return
		}
		utils.RespondInternalError(c, "Failed to fetch sprint", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", sprint)
}

// Update handles PUT /api/sprints/:id
// @Summary Update a sprint
// @Description Update an existing sprint
// @Tags sprints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Param request body request.UpdateSprintRequest true "Update sprint request"
// @Success 200 {object} response.SprintResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id} [put]
func (h *SprintHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	var req request.UpdateSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	sprint, err := h.sprintService.Update(id, &req, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSprintNotFound):
			utils.RespondNotFound(c, "Sprint not found")
		case errors.Is(err, service.ErrInvalidDateRange):
			utils.RespondBadRequest(c, "Invalid date range", err.Error())
		default:
			utils.RespondInternalError(c, "Failed to update sprint", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Sprint updated successfully", sprint)
}

// Delete handles DELETE /api/sprints/:id
// @Summary Delete a sprint
// @Description Delete a sprint by its ID
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id} [delete]
func (h *SprintHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	if err := h.sprintService.Delete(id); err != nil {
		if errors.Is(err, service.ErrSprintNotFound) {
			utils.RespondNotFound(c, "Sprint not found")
			return
		}
		utils.RespondInternalError(c, "Failed to delete sprint", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Sprint deleted successfully", nil)
}

// Start handles POST /api/sprints/:id/start
// @Summary Start a sprint
// @Description Start a sprint (change status from Planning to Active)
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Success 200 {object} response.SprintResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id}/start [post]
func (h *SprintHandler) Start(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	sprint, err := h.sprintService.Start(id, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSprintNotFound):
			utils.RespondNotFound(c, "Sprint not found")
		case errors.Is(err, service.ErrSprintNotPlanning):
			utils.RespondBadRequest(c, "Sprint must be in planning status to start", err.Error())
		case errors.Is(err, service.ErrSprintAlreadyActive):
			utils.RespondError(c, http.StatusConflict, "There is already an active sprint in this project", "SPRINT_ALREADY_ACTIVE", "")
		default:
			utils.RespondInternalError(c, "Failed to start sprint", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Sprint started successfully", sprint)
}

// Complete handles POST /api/sprints/:id/complete
// @Summary Complete a sprint
// @Description Complete a sprint (change status from Active to Completed and calculate velocity)
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Success 200 {object} response.SprintResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id}/complete [post]
func (h *SprintHandler) Complete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	sprint, err := h.sprintService.Complete(id, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSprintNotFound):
			utils.RespondNotFound(c, "Sprint not found")
		case errors.Is(err, service.ErrSprintNotActive):
			utils.RespondBadRequest(c, "Sprint must be in active status to complete", err.Error())
		default:
			utils.RespondInternalError(c, "Failed to complete sprint", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Sprint completed successfully", sprint)
}

// Cancel handles POST /api/sprints/:id/cancel
// @Summary Cancel a sprint
// @Description Cancel a sprint (change status to Cancelled)
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Success 200 {object} response.SprintResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id}/cancel [post]
func (h *SprintHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	sprint, err := h.sprintService.Cancel(id, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSprintNotFound):
			utils.RespondNotFound(c, "Sprint not found")
		case errors.Is(err, service.ErrSprintNotActive):
			utils.RespondBadRequest(c, "Sprint must be in active or planning status to cancel", err.Error())
		default:
			utils.RespondInternalError(c, "Failed to cancel sprint", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Sprint cancelled successfully", sprint)
}

// AddItem handles POST /api/sprints/:id/items
// @Summary Add an item to sprint
// @Description Add a backlog item to a sprint
// @Tags sprints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Param request body request.AddItemToSprintRequest true "Add item request"
// @Success 200 {object} response.SprintWithItemsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id}/items [post]
func (h *SprintHandler) AddItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	var req request.AddItemToSprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request body", err.Error())
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	sprint, err := h.sprintService.AddItem(id, req.ItemID, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSprintNotFound):
			utils.RespondNotFound(c, "Sprint not found")
		case errors.Is(err, service.ErrBacklogItemNotFound):
			utils.RespondNotFound(c, "Backlog item not found")
		case errors.Is(err, service.ErrItemAlreadyInSprint):
			utils.RespondError(c, http.StatusConflict, "Item is already in this sprint", "ITEM_ALREADY_IN_SPRINT", "")
		default:
			utils.RespondInternalError(c, "Failed to add item to sprint", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Item added to sprint successfully", sprint)
}

// RemoveItem handles DELETE /api/sprints/:id/items/:itemId
// @Summary Remove an item from sprint
// @Description Remove a backlog item from a sprint
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Param itemId path string true "Item ID"
// @Success 200 {object} response.SprintWithItemsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id}/items/{itemId} [delete]
func (h *SprintHandler) RemoveItem(c *gin.Context) {
	sprintID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	itemID, err := uuid.Parse(c.Param("itemId"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid item ID", "ID must be a valid UUID")
		return
	}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	sprint, err := h.sprintService.RemoveItem(sprintID, itemID, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSprintNotFound):
			utils.RespondNotFound(c, "Sprint not found")
		case errors.Is(err, service.ErrBacklogItemNotFound):
			utils.RespondNotFound(c, "Backlog item not found")
		case errors.Is(err, service.ErrItemNotInSprint):
			utils.RespondBadRequest(c, "Item is not in this sprint", err.Error())
		default:
			utils.RespondInternalError(c, "Failed to remove item from sprint", err.Error())
		}
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Item removed from sprint successfully", sprint)
}

// GetHistory handles GET /api/sprints/:id/history
// @Summary Get sprint history
// @Description Get the history of a sprint
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Success 200 {array} response.SprintHistoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id}/history [get]
func (h *SprintHandler) GetHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	history, err := h.sprintService.GetHistory(id)
	if err != nil {
		if errors.Is(err, service.ErrSprintNotFound) {
			utils.RespondNotFound(c, "Sprint not found")
			return
		}
		utils.RespondInternalError(c, "Failed to fetch history", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", history)
}

// GetReport handles GET /api/sprints/:id/history/report
// @Summary Get sprint report
// @Description Get the report of a sprint including velocity and completion stats
// @Tags sprints
// @Produce json
// @Security BearerAuth
// @Param id path string true "Sprint ID"
// @Success 200 {object} response.SprintReportResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /sprints/{id}/history/report [get]
func (h *SprintHandler) GetReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondBadRequest(c, "Invalid sprint ID", "ID must be a valid UUID")
		return
	}

	report, err := h.sprintService.GetReport(id)
	if err != nil {
		if errors.Is(err, service.ErrSprintNotFound) {
			utils.RespondNotFound(c, "Sprint not found")
			return
		}
		utils.RespondInternalError(c, "Failed to fetch report", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", report)
}
