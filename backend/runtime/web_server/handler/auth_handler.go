package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authdto "github.com/zcl0621/compx576-smart-dairy-system/dto/auth"
	"github.com/zcl0621/compx576-smart-dairy-system/middleware"
	authservice "github.com/zcl0621/compx576-smart-dairy-system/service/auth"
)

// AuthLogin godoc
// @Summary login
// @Description login with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdto.LoginRequest true "login request"
// @Success 200 {object} authdto.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *Handler) AuthLogin(c *gin.Context) {
	var request authdto.LoginRequest
	if !bindJSON(c, &request) {
		return
	}

	response, err := authservice.Login(&request)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// AuthRefreshToken godoc
// @Summary refresh token
// @Description get new token for current user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} authdto.LoginResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/refresh [post]
func (h *Handler) AuthRefreshToken(c *gin.Context) {
	userID, err := middleware.GetAuthUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	response, err := authservice.RefreshToken(userID)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// AuthPasswordResetRequest godoc
// @Summary request password reset code
// @Description send reset code to email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdto.PasswordResetRequestEmailRequest true "reset email request"
// @Success 200 {object} authdto.MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/password-reset/request [post]
func (h *Handler) AuthPasswordResetRequest(c *gin.Context) {
	var request authdto.PasswordResetRequestEmailRequest
	if !bindJSON(c, &request) {
		return
	}

	response, err := authservice.RequestPasswordReset(&request)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// AuthPasswordResetVerify godoc
// @Summary verify password reset code
// @Description check reset code and return reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdto.PasswordResetVerifyRequest true "reset verify request"
// @Success 200 {object} authdto.PasswordResetVerifyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/password-reset/verify [post]
func (h *Handler) AuthPasswordResetVerify(c *gin.Context) {
	var request authdto.PasswordResetVerifyRequest
	if !bindJSON(c, &request) {
		return
	}

	response, err := authservice.VerifyPasswordResetCode(&request)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// AuthPasswordResetConfirm godoc
// @Summary confirm password reset
// @Description reset password with reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdto.PasswordResetConfirmRequest true "reset confirm request"
// @Success 200 {object} authdto.MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/auth/password-reset/confirm [post]
func (h *Handler) AuthPasswordResetConfirm(c *gin.Context) {
	var request authdto.PasswordResetConfirmRequest
	if !bindJSON(c, &request) {
		return
	}

	response, err := authservice.ConfirmPasswordReset(&request)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
