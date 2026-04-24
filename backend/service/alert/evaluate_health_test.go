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

func TestEvaluateHealth_AllCritical_CreatesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HealthCow1", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.0, now.Add(-60*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.2, now.Add(-30*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.5, now)

		err := alert.EvaluateHealth(c.ID, model.MetricTypeTemperature)

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND severity = 'critical' AND status = 'active'", c.ID, model.MetricTypeTemperature).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestEvaluateHealth_AllWarning_CreatesWarningAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HealthCow2", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 40.0, now.Add(-60*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 40.1, now.Add(-30*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 40.2, now)

		err := alert.EvaluateHealth(c.ID, model.MetricTypeTemperature)

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND severity = 'warning' AND status = 'active'", c.ID, model.MetricTypeTemperature).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestEvaluateHealth_LatestNormal_ResolvesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HealthCow3", model.CowStatusInFarm)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.0, now.Add(-60*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 40.0, now.Add(-30*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.8, now) // normal

		err := alert.EvaluateHealth(c.ID, model.MetricTypeTemperature)

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeTemperature).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestEvaluateHealth_FewerThan3_SkipsCreate(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HealthCow4", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.0, now.Add(-30*time.Second))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 41.2, now)

		err := alert.EvaluateHealth(c.ID, model.MetricTypeTemperature)

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND status = 'active'", c.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestEvaluateHealth_FewerThan3_StillResolves(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HealthCow5", model.CowStatusInFarm)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityCritical, model.AlertStatusActive)
		now := time.Now()
		// only 1 reading but it's normal — resolve path runs independently
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.5, now)

		err := alert.EvaluateHealth(c.ID, model.MetricTypeTemperature)

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeTemperature).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
