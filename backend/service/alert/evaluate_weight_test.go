package alert_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/alert"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestEvaluateWeight_DeviationAboveThreshold_CreatesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "WeightCow1", model.CowStatusInFarm)
		now := time.Now()
		// avg = 500kg; 349kg → deviation = |349-500|/500 = 30.2% → above 30% threshold
		for i := 1; i <= 7; i++ {
			testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeWeight, 500.0, now.AddDate(0, 0, -i))
		}

		err := alert.EvaluateWeight(c.ID, 349.0, now.Add(-time.Minute))

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeWeight).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestEvaluateWeight_DeviationBelowResolveThreshold_Resolves(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "WeightCow2", model.CowStatusInFarm)
		now := time.Now()
		for i := 1; i <= 7; i++ {
			testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeWeight, 500.0, now.AddDate(0, 0, -i))
		}
		// seed an active weight alert
		a := testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		tx.Model(a).Update("metric_key", model.MetricTypeWeight)

		// 498kg → deviation = 0.4% → below 15% resolve threshold
		err := alert.EvaluateWeight(c.ID, 498.0, now.Add(-time.Minute))

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeWeight).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestEvaluateWeight_NotEnoughHistory_Skips(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "WeightCow3", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeWeight, 500.0, now.AddDate(0, 0, -1))

		err := alert.EvaluateWeight(c.ID, 100.0, now.Add(-time.Minute))

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND status = 'active'", c.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
