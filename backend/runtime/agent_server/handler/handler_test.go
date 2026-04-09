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
	"github.com/zcl0621/compx576-smart-dairy-system/middleware"
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
	r.POST("/api/token", h.Token)
	return r
}

func setupMetricRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewHandler()
	r.POST("/api/token", h.Token)
	r.POST("/api/metric", middleware.NeedCowAuth(), h.Metric)
	return r
}

func getToken(t *testing.T, router *gin.Engine, cowID string) string {
	t.Helper()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/token?cow_id="+cowID, nil)
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	return resp["token"].(string)
}

func TestToken_ValidCow(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Daisy", model.CowStatusInFarm, model.CowConditionNormal)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/token?cow_id="+cow.Tag, nil)
		setupRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.NotEmpty(t, resp["token"])
	})
}

func TestToken_AnyCowID(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		// token endpoint issues token for any cow_id without db lookup
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/token?cow_id=any-tag", nil)
		setupRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestToken_MissingCowID(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/token", nil)
		setupRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMetric_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)
	require.NoError(t, mq.Init())

	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Bella", model.CowStatusInFarm, model.CowConditionNormal)
		router := setupMetricRouter()
		token := getToken(t, router, cow.Tag)

		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       cow.Tag,
			"source":       "cow_agent",
			"metric_type":  "temperature",
			"metric_value": 38.5,
			"unit":         "celsius",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
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
	setupMetricRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMetric_CowIDMismatch(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)
	require.NoError(t, mq.Init())

	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Clara", model.CowStatusInFarm, model.CowConditionNormal)
		router := setupMetricRouter()
		token := getToken(t, router, cow.Tag)

		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       "different-cow-tag",
			"source":       "cow_agent",
			"metric_type":  "temperature",
			"metric_value": 38.5,
			"unit":         "celsius",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestMetric_InvalidSource(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)
	require.NoError(t, mq.Init())

	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Dana", model.CowStatusInFarm, model.CowConditionNormal)
		router := setupMetricRouter()
		token := getToken(t, router, cow.Tag)

		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       cow.Tag,
			"source":       "invalid_source",
			"metric_type":  "temperature",
			"metric_value": 38.5,
			"unit":         "celsius",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMetric_InvalidMetricType(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)
	require.NoError(t, mq.Init())

	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "Eva", model.CowStatusInFarm, model.CowConditionNormal)
		router := setupMetricRouter()
		token := getToken(t, router, cow.Tag)

		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       cow.Tag,
			"source":       "cow_agent",
			"metric_type":  "invalid_type",
			"metric_value": 38.5,
			"unit":         "celsius",
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
