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

func TestEvaluateOffline_NeverSentMetric_CreatesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "OfflineCow1", model.CowStatusInFarm)

		err := alert.EvaluateOffline()

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeDevice).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestEvaluateOffline_MetricOlderThan10Min_CreatesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "OfflineCow2", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.5, time.Now().Add(-15*time.Minute))

		err := alert.EvaluateOffline()

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeDevice).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestEvaluateOffline_RecentMetric_ResolvesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "OfflineCow3", model.CowStatusInFarm)
		a := testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityOffline, model.AlertStatusActive)
		tx.Model(a).Update("metric_key", model.MetricTypeDevice)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.5, time.Now().Add(-2*time.Minute))

		err := alert.EvaluateOffline()

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeDevice).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestEvaluateOffline_SoldCow_Ignored(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "SoldCow", model.CowStatusSold)

		err := alert.EvaluateOffline()

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND status = 'active'", c.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
