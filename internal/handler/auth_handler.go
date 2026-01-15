package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sprint-backlog/internal/dto/request"
	"sprint-backlog/internal/service"
	"sprint-backlog/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// VerifyGoogleCode handles POST /api/auth/google/verify
// @Summary Verify Google OAuth code
// @Description Verify Google OAuth authorization code and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.GoogleVerifyRequest true "Google verify request"
// @Success 200 {object} response.AuthResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/google/verify [post]
func (h *AuthHandler) VerifyGoogleCode(c *gin.Context) {
	var req request.GoogleVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, "Invalid request", err.Error())
		return
	}

	authResp, err := h.authService.VerifyGoogleCode(c.Request.Context(), req.Code, req.RedirectURI)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Authentication failed", "AUTH_FAILED", err.Error())
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "Authentication successful", authResp)
}

// GetCurrentUser handles GET /api/auth/me
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDValue, exists := c.Get("user_id")
	if !exists {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		utils.RespondUnauthorized(c, "Invalid user ID")
		return
	}

	userResp, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		utils.RespondNotFound(c, "User not found")
		return
	}

	utils.RespondSuccess(c, http.StatusOK, "", userResp)
}
