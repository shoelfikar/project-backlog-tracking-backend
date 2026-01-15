package request

// GoogleVerifyRequest represents the request to verify Google OAuth code
type GoogleVerifyRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirect_uri" binding:"required"`
}
