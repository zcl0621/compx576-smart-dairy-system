package handler

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/middleware"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	"go.uber.org/zap"
)

// tag to cow_id cache so we don't hit db on every metric
var (
	tagCache   = map[string]string{}
	tagCacheMu sync.RWMutex
)

func resolveCowID(tag string) (string, bool) {
	tagCacheMu.RLock()
	id, ok := tagCache[tag]
	tagCacheMu.RUnlock()
	if ok {
		return id, true
	}

	var cow model.Cow
	if err := pg.DB.Where("tag = ?", tag).First(&cow).Error; err != nil {
		return "", false
	}

	tagCacheMu.Lock()
	tagCache[tag] = cow.ID
	tagCacheMu.Unlock()
	return cow.ID, true
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// Token generates a cow JWT for device auth
func (h *Handler) Token(c *gin.Context) {
	cowID := c.Query("cow_id")
	if cowID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cow_id is required"})
		return
	}

	token, expiresAt, err := middleware.GenerateCowToken(cowID)
	if err != nil {
		projectlog.L().Error("generate cow token failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": expiresAt.Unix(),
	})
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

	// check cow_id matches JWT
	tokenCowID, err := middleware.GetAuthCowID(c)
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

	// agent sends cow tag, resolve to cow_id for db storage
	cowID, ok := resolveCowID(req.CowID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "cow not found"})
		return
	}

	fields := map[string]string{
		"cow_id":       cowID,
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
