package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	alertdto "github.com/zcl0621/compx576-smart-dairy-system/dto/alert"
	alertservice "github.com/zcl0621/compx576-smart-dairy-system/service/alert"
)

// AlertSummary godoc
// @Summary get alert summary
// @Description get active alert count
// @Tags Alert
// @Produce json
// @Security BearerAuth
// @Success 200 {object} alertdto.SummaryResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/alert/summary [get]
func (h *Handler) AlertSummary(c *gin.Context) {
	response, err := alertservice.SummaryService()
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// AlertList godoc
// @Summary list alerts
// @Description get alert list
// @Tags Alert
// @Produce json
// @Security BearerAuth
// @Param page query int false "page num" default(1)
// @Param page_size query int false "page size" default(20)
// @Param cow_id query string false "Cow ID"
// @Param severity query string false "alert level"
// @Success 200 {object} alertdto.ListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/alert/list [get]
func (h *Handler) AlertList(c *gin.Context) {
	var query alertdto.ListQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := alertservice.ListService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
