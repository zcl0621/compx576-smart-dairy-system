package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zcl0621/compx576-smart-dairy-system/middleware"
	authservice "github.com/zcl0621/compx576-smart-dairy-system/service/auth"
	metricservice "github.com/zcl0621/compx576-smart-dairy-system/service/metric"
	"gorm.io/gorm"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Health godoc
// @Summary health check
// @Description check if api is up
// @Tags Health
// @Success 200 {object} OKResponse
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func bindQuery(c *gin.Context, target any) bool {
	if err := c.ShouldBindQuery(target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}

func bindJSON(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}

func writeOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func writeError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	}
	if errors.Is(err, metricservice.ErrBadMetricRange) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, authservice.ErrBadLogin) {
		status = http.StatusUnauthorized
	}
	if errors.Is(err, authservice.ErrResetCodeNotFound) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, authservice.ErrResetTokenNotFound) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, middleware.ErrBadToken) {
		status = http.StatusUnauthorized
	}

	c.JSON(status, gin.H{"error": err.Error()})
}
