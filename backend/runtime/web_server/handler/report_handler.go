package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	reportdto "github.com/zcl0621/compx576-smart-dairy-system/dto/report"
	reportservice "github.com/zcl0621/compx576-smart-dairy-system/service/report"
)

// ReportList godoc
// @Summary list reports
// @Description get report list
// @Tags Report
// @Produce json
// @Security BearerAuth
// @Param page query int false "page num" default(1)
// @Param page_size query int false "page size" default(20)
// @Success 200 {object} reportdto.ListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/report/list [get]
func (h *Handler) ReportList(c *gin.Context) {
	var query reportdto.ListQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := reportservice.ListService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ReportLatest godoc
// @Summary get latest report
// @Description get latest report for one cow
// @Tags Report
// @Produce json
// @Security BearerAuth
// @Param cow_id query string true "Cow ID"
// @Success 200 {object} reportdto.ReportItem
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/report/latest [get]
func (h *Handler) ReportLatest(c *gin.Context) {
	var query reportdto.LatestQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := reportservice.LatestService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
