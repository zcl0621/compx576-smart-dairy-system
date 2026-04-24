package agent_server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/mq"
	agent_server "github.com/zcl0621/compx576-smart-dairy-system/runtime/agent_server"
	"github.com/zcl0621/compx576-smart-dairy-system/runtime/agent_server/consumer"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestIntegration_PostMetricEndToEnd(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.FlushRedis(t)

	testhelper.WithTx(t, func(tx *gorm.DB) {
		// mq.Init must be called INSIDE WithTx because WithTx calls FlushRedis
		// at start which deletes the stream and consumer group
		require.NoError(t, mq.Init())

		cow := testhelper.SeedCow(t, tx, "Fiona", model.CowStatusInFarm)

		router := agent_server.NewRouter()

		// 1. get token using tag (like real agent)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/token?cow_id="+cow.Tag, nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var tokenResp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tokenResp))
		token := tokenResp["token"].(string)

		// 2. start consumer
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			consumer.StartMetricWriter(ctx)
			close(done)
		}()

		// 3. post metric using tag (agent server resolves to cow_id)
		body, _ := json.Marshal(map[string]interface{}{
			"cow_id":       cow.Tag,
			"source":       "cow_agent",
			"metric_type":  "heart_rate",
			"metric_value": 75.0,
			"unit":         "bpm",
		})

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/api/metric", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// 4. wait for consumer to write
		time.Sleep(2 * time.Second)
		cancel()
		<-done

		// 5. verify in database
		var metric model.Metric
		err := pg.DB.Where("cow_id = ? AND metric_type = ?", cow.ID, "heart_rate").First(&metric).Error
		require.NoError(t, err)
		assert.Equal(t, cow.ID, metric.CowID)
		assert.InDelta(t, 75.0, metric.MetricValue, 0.01)
		assert.Equal(t, model.MetricUnitBPM, metric.Unit)
	})
}
