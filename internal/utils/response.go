package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ErrUserNotInContext is returned when user_id is not found in context
var ErrUserNotInContext = errors.New("user not found in context")

// GetUserIDFromContext extracts user ID from gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, ErrUserNotInContext
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrUserNotInContext
	}

	return userID, nil
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorDetail represents error details
type ErrorDetail struct {
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   ErrorDetail `json:"error,omitempty"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// RespondSuccess sends a success response
func RespondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondError sends an error response
func RespondError(c *gin.Context, statusCode int, message string, code string, details string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Error: ErrorDetail{
			Code:    code,
			Details: details,
		},
	})
}

// RespondBadRequest sends a 400 Bad Request response
func RespondBadRequest(c *gin.Context, message string, details string) {
	RespondError(c, http.StatusBadRequest, message, "BAD_REQUEST", details)
}

// RespondUnauthorized sends a 401 Unauthorized response
func RespondUnauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, message, "UNAUTHORIZED", "")
}

// RespondForbidden sends a 403 Forbidden response
func RespondForbidden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, message, "FORBIDDEN", "")
}

// RespondNotFound sends a 404 Not Found response
func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, message, "NOT_FOUND", "")
}

// RespondInternalError sends a 500 Internal Server Error response
func RespondInternalError(c *gin.Context, message string, details ...string) {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}
	RespondError(c, http.StatusInternalServerError, message, "INTERNAL_ERROR", detail)
}

// RespondPaginated sends a paginated response
func RespondPaginated(c *gin.Context, data interface{}, page, perPage, total int) {
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
