package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sprint-backlog/internal/service"
	"sprint-backlog/internal/utils"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetAll returns all users
// @Summary Get all users
// @Description Get all registered users
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse{data=[]response.UserResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll()
	if err != nil {
		utils.RespondInternalError(c, "Failed to get users", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Users retrieved successfully", users)
}

// GetByID returns a user by ID
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} utils.SuccessResponse{data=response.UserResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid user ID", err.Error())
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		utils.RespondInternalError(c, "Failed to get user", err.Error())
		return
	}
	if user == nil {
		utils.RespondNotFound(c, "User not found")
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "User retrieved successfully", user)
}

// GetActivities returns a user's activity history
// @Summary Get user activities
// @Description Get a user's activity history (item and sprint actions)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param limit query int false "Maximum number of activities to return" default(50)
// @Success 200 {object} utils.SuccessResponse{data=response.UserActivitiesResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /users/{id}/activities [get]
func (h *UserHandler) GetActivities(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondBadRequest(c, "Invalid user ID", err.Error())
		return
	}

	// Get limit from query params (default 50)
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	activities, err := h.userService.GetActivities(id, limit)
	if err != nil {
		utils.RespondInternalError(c, "Failed to get user activities", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "User activities retrieved successfully", activities)
}

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Description Update the current user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} utils.SuccessResponse{data=response.UserResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 501 {object} utils.ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	utils.RespondError(c, http.StatusNotImplemented, "Not implemented yet", "NOT_IMPLEMENTED", "")
}
