package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	"go.uber.org/zap"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

type metricRequest struct {
	CowID       string   `json:"cow_id" binding:"required"`
	Source      string   `json:"source" binding:"required"`
	MetricType  string   `json:"metric_type" binding:"required"`
	MetricValue *float64 `json:"metric_value" binding:"required"`
	Unit        string   `json:"unit" binding:"required"`
}

var validSources = map[string]bool{
	string(model.MetricSourceCowAgent):       true,
	string(model.MetricSourceMilkingMachine): true,
	string(model.MetricSourceWeightMachine):  true,
}

var validUnits = map[string]bool{
	string(model.MetricUnitCelsius): true,
	string(model.MetricUnitBPM):     true,
	string(model.MetricUnitPercent): true,
	string(model.MetricUnitDegrees): true,
	string(model.MetricUnitLiters):  true,
	string(model.MetricUnitSeconds): true,
	string(model.MetricUnitKG):      true,
}

// validMetricTypes excludes device — that type is reserved for device status tracking, not agent data
var validMetricTypes = map[string]bool{
	string(model.MetricTypeTemperature):     true,
	string(model.MetricTypeHeartRate):       true,
	string(model.MetricTypeBloodOxygen):     true,
	string(model.MetricTypeLatitude):        true,
	string(model.MetricTypeLongitude):       true,
	string(model.MetricTypeMilkAmount):      true,
	string(model.MetricTypeMilkingDuration): true,
	string(model.MetricTypeWeight):          true,
}

// Metric receives a single metric and pushes to Redis Stream
func (h *Handler) Metric(c *gin.Context) {
	var req metricRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := bearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "bad token"})
		return
	}

	tokenCowID, err := resolveCowIDForToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "bad token"})
		return
	}
	if tokenCowID != req.CowID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cow_id mismatch"})
		return
	}

	if !validSources[req.Source] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid source"})
		return
	}
	if !validMetricTypes[req.MetricType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid metric_type"})
		return
	}
	if !validUnits[req.Unit] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid unit"})
		return
	}

	fields := map[string]string{
		"cow_id":       req.CowID,
		"source":       req.Source,
		"metric_type":  req.MetricType,
		"metric_value": strconv.FormatFloat(*req.MetricValue, 'f', 2, 64),
		"unit":         req.Unit,
		"timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
	}

	if err := mq.Publish(fields); err != nil {
		projectlog.L().Error("publish metric failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func resolveCowIDForToken(token string) (string, error) {
	key := "agent_token:" + token
	if cached, err := redis.Get(key); err == nil && cached != "" {
		return cached, nil
	}

	var cow model.Cow
	if err := pg.DB.Where("agent_token = ?", token).First(&cow).Error; err != nil {
		return "", err
	}

	_ = redis.Set(key, cow.ID, 24*time.Hour)
	return cow.ID, nil
}

func bearerToken(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] == "" {
		return "", errors.New("bad token")
	}
	return parts[1], nil
}
