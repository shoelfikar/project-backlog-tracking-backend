package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sprint-backlog/internal/dto/response"
	"sprint-backlog/internal/models"
	"sprint-backlog/internal/repository"
	"sprint-backlog/internal/utils"
)

type AuthService interface {
	VerifyGoogleCode(ctx context.Context, code, redirectURI string) (*response.AuthResponse, error)
	GetCurrentUser(userID uuid.UUID) (*response.UserResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// VerifyGoogleCode verifies the Google OAuth code and returns JWT + user data
func (s *authService) VerifyGoogleCode(ctx context.Context, code, redirectURI string) (*response.AuthResponse, error) {
	// 1. Exchange code for tokens
	tokenResp, err := utils.ExchangeCodeForToken(ctx, code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// 2. Decode ID token to get user info
	googleUser, err := utils.DecodeIDToken(tokenResp.IDToken)
	if err != nil {
		// Fallback: try to get user info from access token
		googleUser, err = utils.GetGoogleUserInfo(ctx, tokenResp.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to get user info: %w", err)
		}
	}

	// 3. Find or create user
	user, err := s.findOrCreateUser(googleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to find/create user: %w", err)
	}

	// 4. Generate JWT token (expiry follows Google token expiry)
	token, err := utils.GenerateToken(user.ID, user.Email, user.GoogleID, tokenResp.ExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 5. Return auth response
	return &response.AuthResponse{
		Token:     token,
		ExpiresIn: tokenResp.ExpiresIn,
		User: response.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
		},
	}, nil
}

// findOrCreateUser finds an existing user or creates a new one
func (s *authService) findOrCreateUser(googleUser *utils.GoogleUserInfo) (*models.User, error) {
	// Try to find by Google ID first
	user, err := s.userRepo.GetByGoogleID(googleUser.ID)
	if err != nil {
		return nil, err
	}

	if user != nil {
		// User exists, update info if changed
		updated := false
		if user.Name != googleUser.Name {
			user.Name = googleUser.Name
			updated = true
		}
		if googleUser.Picture != "" && (user.AvatarURL == nil || *user.AvatarURL != googleUser.Picture) {
			user.AvatarURL = &googleUser.Picture
			updated = true
		}

		if updated {
			if err := s.userRepo.Update(user); err != nil {
				return nil, err
			}
		}

		return user, nil
	}

	// Try to find by email (user might have been created differently)
	user, err = s.userRepo.GetByEmail(googleUser.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		// Link Google ID to existing user
		user.GoogleID = googleUser.ID
		if googleUser.Picture != "" {
			user.AvatarURL = &googleUser.Picture
		}
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// Create new user
	var avatarURL *string
	if googleUser.Picture != "" {
		avatarURL = &googleUser.Picture
	}

	newUser := &models.User{
		GoogleID:  googleUser.ID,
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		AvatarURL: avatarURL,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// GetCurrentUser returns the current user's info
func (s *authService) GetCurrentUser(userID uuid.UUID) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return &response.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
	}, nil
}
