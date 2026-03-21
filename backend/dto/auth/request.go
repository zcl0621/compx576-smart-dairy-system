package auth

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type PasswordResetRequestEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetVerifyRequest struct {
	Code string `json:"code" binding:"required,len=6,numeric"`
}

type PasswordResetConfirmRequest struct {
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
