package auth

import "time"

type LoginUser struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      LoginUser `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type PasswordResetVerifyResponse struct {
	ResetToken string `json:"reset_token"`
}
