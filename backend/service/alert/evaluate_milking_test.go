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

func TestEvaluateMilking_BelowThreshold_CreatesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MilkCow1", model.CowStatusInFarm)
		now := time.Now()
		// 7 days of history at 10L each
		for i := 1; i <= 7; i++ {
			testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeMilkAmount, 10.0, now.AddDate(0, 0, -i))
		}
		// current: 6L = 60% of 10L avg → below 70% threshold
		err := alert.EvaluateMilking(c.ID, 6.0, now.Add(-time.Minute))

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeMilkAmount).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestEvaluateMilking_AboveResolveThreshold_Resolves(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MilkCow2", model.CowStatusInFarm)
		now := time.Now()
		for i := 1; i <= 7; i++ {
			testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeMilkAmount, 10.0, now.AddDate(0, 0, -i))
		}
		// seed an active milk_amount alert to resolve
		a := testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		tx.Model(a).Update("metric_key", model.MetricTypeMilkAmount)

		// current: 9L = 90% of avg → above 85% resolve threshold
		err := alert.EvaluateMilking(c.ID, 9.0, now.Add(-time.Minute))

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeMilkAmount).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestEvaluateMilking_NotEnoughHistory_Skips(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MilkCow3", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeMilkAmount, 10.0, now.AddDate(0, 0, -1))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeMilkAmount, 10.0, now.AddDate(0, 0, -2))

		err := alert.EvaluateMilking(c.ID, 1.0, now.Add(-time.Minute))

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND status = 'active'", c.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
