package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cowdto "github.com/zcl0621/compx576-smart-dairy-system/dto/cow"
	metricdto "github.com/zcl0621/compx576-smart-dairy-system/dto/metric"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	cowservice "github.com/zcl0621/compx576-smart-dairy-system/service/cow"
	metricservice "github.com/zcl0621/compx576-smart-dairy-system/service/metric"
)

// CowList godoc
// @Summary list cows
// @Description get cow list
// @Tags Cow
// @Produce json
// @Security BearerAuth
// @Param page query int false "page num" default(1)
// @Param page_size query int false "page size" default(20)
// @Param name query string false "cow name"
// @Param condition query string false "cow condition"
// @Param status query string false "cow status"
// @Param sort query string false "sort field" default(updated_at)
// @Success 200 {object} cowdto.ListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/list [get]
func (h *Handler) CowList(c *gin.Context) {
	var query cowdto.ListQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := cowservice.CowListService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CowInfo godoc
// @Summary get cow info
// @Description get one cow by id
// @Tags Cow
// @Produce json
// @Security BearerAuth
// @Param id query string true "Cow ID"
// @Success 200 {object} cowdto.InfoResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/info [get]
func (h *Handler) CowInfo(c *gin.Context) {
	var query cowdto.InfoQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := cowservice.CowInfoService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CowCreate godoc
// @Summary create cow
// @Description create one cow
// @Tags Cow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body cowdto.CreateRequest true "create cow request"
// @Success 200 {object} OKResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/create [post]
func (h *Handler) CowCreate(c *gin.Context) {
	var request cowdto.CreateRequest
	if !bindJSON(c, &request) {
		return
	}

	if err := cowservice.CowCreateService(&request); err != nil {
		writeError(c, err)
		return
	}
	writeOK(c)
}

// CowUpdate godoc
// @Summary update cow
// @Description update one cow
// @Tags Cow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body cowdto.UpdateRequest true "update cow request"
// @Success 200 {object} OKResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/update [post]
func (h *Handler) CowUpdate(c *gin.Context) {
	var request cowdto.UpdateRequest
	if !bindJSON(c, &request) {
		return
	}

	if err := cowservice.CowUpdateService(&request); err != nil {
		writeError(c, err)
		return
	}
	writeOK(c)
}

// CowMetricTemperature godoc
// @Summary get temperature metric
// @Description get temperature data for one cow
// @Tags Cow Metric
// @Produce json
// @Security BearerAuth
// @Param cow_id query string true "Cow ID"
// @Param range query string false "range" Enums(24h,7d,30d,all) default(24h)
// @Success 200 {object} cowdto.TemperatureMetricResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/metric/temperature [get]
func (h *Handler) CowMetricTemperature(c *gin.Context) {
	var query cowdto.MetricQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := metricservice.TemperatureService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CowMetricHeartRate godoc
// @Summary get heart rate metric
// @Description get heart rate data for one cow
// @Tags Cow Metric
// @Produce json
// @Security BearerAuth
// @Param cow_id query string true "Cow ID"
// @Param range query string false "range" Enums(24h,7d,30d,all) default(24h)
// @Success 200 {object} cowdto.HeartRateMetricResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/metric/heart_rate [get]
func (h *Handler) CowMetricHeartRate(c *gin.Context) {
	var query cowdto.MetricQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := metricservice.HeartRateService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CowMetricBloodOxygen godoc
// @Summary get blood oxygen metric
// @Description get blood oxygen data for one cow
// @Tags Cow Metric
// @Produce json
// @Security BearerAuth
// @Param cow_id query string true "Cow ID"
// @Param range query string false "range" Enums(24h,7d,30d,all) default(24h)
// @Success 200 {object} cowdto.BloodOxygenMetricResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/metric/blood_oxygen [get]
func (h *Handler) CowMetricBloodOxygen(c *gin.Context) {
	var query cowdto.MetricQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := metricservice.BloodOxygenService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CowMetricMilkAmount godoc
// @Summary get milk amount metric
// @Description get milk amount data for one cow
// @Tags Cow Metric
// @Produce json
// @Security BearerAuth
// @Param cow_id query string true "Cow ID"
// @Param range query string false "range" Enums(24h,7d,30d,all) default(24h)
// @Success 200 {object} cowdto.MilkAmountMetricResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/metric/milk_amount [get]
func (h *Handler) CowMetricMilkAmount(c *gin.Context) {
	var query cowdto.MetricQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := metricservice.MilkAmountService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CowMetricMovement godoc
// @Summary get movement metric
// @Description get movement data for one cow
// @Tags Cow Metric
// @Produce json
// @Security BearerAuth
// @Param cow_id query string true "Cow ID"
// @Param range query string false "range" Enums(24h,7d,30d,all) default(24h)
// @Success 200 {object} cowdto.MovementMetricResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/metric/movement [get]
func (h *Handler) CowMetricMovement(c *gin.Context) {
	var query cowdto.MetricQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := metricservice.MovementService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CowMetricWeight godoc
// @Summary get weight metric
// @Description get weight data for one cow
// @Tags Cow Metric
// @Produce json
// @Security BearerAuth
// @Param cow_id query string true "Cow ID"
// @Param range query string false "range" Enums(24h,7d,30d,all) default(24h)
// @Success 200 {object} cowdto.WeightMetricResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cow/metric/weight [get]
func (h *Handler) CowMetricWeight(c *gin.Context) {
	var query cowdto.MetricQuery
	if !bindQuery(c, &query) {
		return
	}

	response, err := metricservice.WeightService(&query)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// MetricList godoc
// @Summary list raw metric records
// @Description get paginated metric records with optional filtering
// @Tags Metric
// @Produce json
// @Security BearerAuth
// @Param page query int false "page num" default(1)
// @Param page_size query int false "page size" default(20)
// @Param cow_id query string false "cow id"
// @Param metric_type query string false "metric type"
// @Success 200 {object} metricdto.ListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/metric/list [get]
func (h *Handler) MetricList(c *gin.Context) {
	var query metricdto.ListQuery
	if !bindQuery(c, &query) {
		return
	}
	resp, err := metricservice.MetricListService(&query)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CowMetricMovementPath godoc
// @Summary get cow movement path for map
// @Description return GPS path points with stay detection for map rendering
// @Tags Cow
// @Security BearerAuth
// @Param cow_id query string true "cow id"
// @Param range query string false "range" Enums(24h, 7d, 30d, all) default(24h)
// @Success 200 {object} cowdto.MovementPathResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/cow/metric/movement_path [get]
func (h *Handler) CowMetricMovementPath(c *gin.Context) {
	var req cowdto.MetricQuery
	if !bindQuery(c, &req) {
		return
	}

	resp, err := metricservice.MovementPathService(&metricservice.MetricQuery{
		CowID:       req.CowID,
		MetricRange: model.MetricRange(req.Range),
	})
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
