package request

// UpdateProfileRequest represents the request body for updating user profile
type UpdateProfileRequest struct {
	Name      string `json:"name" binding:"omitempty,min=2,max=100"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}
