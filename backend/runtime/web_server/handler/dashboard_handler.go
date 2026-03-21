package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dashboarddto "github.com/zcl0621/compx576-smart-dairy-system/dto/dashboard"
	dashboardservice "github.com/zcl0621/compx576-smart-dairy-system/service/dashboard"
)

// DashboardSummary godoc
// @Summary get dashboard summary
// @Description get dashboard cow count
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dashboarddto.SummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/dashboard/summary [get]
func (h *Handler) DashboardSummary(c *gin.Context) {
	response, err := dashboardservice.SummaryService()
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DashboardList godoc
// @Summary list dashboard cows
// @Description get dashboard cow list
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Param page query int false "page num" default(1)
// @Param page_size query int false "page size" default(20)
// @Success 200 {object} dashboarddto.ListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/dashboard/list [get]
func (h *Handler) DashboardList(c *gin.Context) {
	var query dashboarddto.ListQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := dashboardservice.ListService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
