package consumer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/runtime/agent_server/consumer"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestHandleAlertEval_Temperature_CallsHealthEvaluator(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "AlertConsumerCow1", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.0, now.Add(-60*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.2, now.Add(-30*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.5, now.Add(-time.Second))

		err := consumer.HandleAlertEval("msg-1", map[string]any{
			"cow_id":       c.ID,
			"source":       "cow_agent",
			"metric_type":  "temperature",
			"metric_value": "41.5",
			"unit":         "celsius",
			"timestamp":    "0",
		})

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeTemperature).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestHandleAlertEval_MilkingDuration_Skipped(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "AlertConsumerCow2", model.CowStatusInFarm)

		err := consumer.HandleAlertEval("msg-2", map[string]any{
			"cow_id":       c.ID,
			"source":       "milking_machine",
			"metric_type":  "milking_duration",
			"metric_value": "400",
			"unit":         "seconds",
			"timestamp":    "0",
		})

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND status = 'active'", c.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestHandleAlertEval_MissingCowID_Skipped(t *testing.T) {
	err := consumer.HandleAlertEval("msg-3", map[string]any{
		"metric_type":  "temperature",
		"metric_value": "41.0",
	})
	require.NoError(t, err)
}

func TestHandleAlertEval_BadMetricValue_Skipped(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "AlertConsumerCow3", model.CowStatusInFarm)

		err := consumer.HandleAlertEval("msg-4", map[string]any{
			"cow_id":       c.ID,
			"metric_type":  "milk_amount",
			"metric_value": "not-a-number",
			"timestamp":    "0",
		})

		require.NoError(t, err)
	})
}
