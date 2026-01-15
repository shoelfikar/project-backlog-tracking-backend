package response

import (
	"github.com/google/uuid"

	"sprint-backlog/internal/models"
)

// UserResponse represents the user data in responses
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL *string   `json:"avatar_url"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresIn int          `json:"expires_in"`
	User      UserResponse `json:"user"`
}

// ToUserResponse converts a User model to UserResponse
func ToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}
}
