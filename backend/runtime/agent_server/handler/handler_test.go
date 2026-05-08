package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	"github.com/zcl0621/compx576-smart-dairy-system/runtime/agent_server/handler"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewHandler()
	r.POST("/api/metric", h.Metric)
	return r
}

func TestMetric_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)
	require.NoError(t, mq.Init())

	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Bella", model.CowStatusInFarm)
		router := setupRouter()

		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       cow.ID,
			"source":       "cow_agent",
			"metric_type":  "temperature",
			"metric_value": 38.5,
			"unit":         "celsius",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+cow.AgentToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		length, err := mq.StreamLen()
		require.NoError(t, err)
		assert.Equal(t, int64(1), length)
	})
}

func TestMetric_NoAuth(t *testing.T) {
	testhelper.SetupTestDB(t)

	body, _ := json.Marshal(map[string]interface{}{
		"cow_id":       "cow-1",
		"source":       "cow_agent",
		"metric_type":  "temperature",
		"metric_value": 38.5,
		"unit":         "celsius",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	setupRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMetric_CowIDMismatch(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)
	require.NoError(t, mq.Init())

	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Clara", model.CowStatusInFarm)
		router := setupRouter()

		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       "different-cow-id",
			"source":       "cow_agent",
			"metric_type":  "temperature",
			"metric_value": 38.5,
			"unit":         "celsius",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+cow.AgentToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestMetric_InvalidSource(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)
	require.NoError(t, mq.Init())

	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Dana", model.CowStatusInFarm)
		router := setupRouter()

		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       cow.ID,
			"source":       "invalid_source",
			"metric_type":  "temperature",
			"metric_value": 38.5,
			"unit":         "celsius",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+cow.AgentToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMetric_InvalidMetricType(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)
	require.NoError(t, mq.Init())

	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Eva", model.CowStatusInFarm)
		router := setupRouter()

		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       cow.ID,
			"source":       "cow_agent",
			"metric_type":  "invalid_type",
			"metric_value": 38.5,
			"unit":         "celsius",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+cow.AgentToken)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
