package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=18"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=18"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"` // "user" или "admin"
	ExpiresIn   int    `json:"expires_in"` // секунды
}

type MeResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"` // "user" или "admin"
}
