package alert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/alert"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestCreateIfNotExists_CreatesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Cow1", model.CowStatusInFarm)

		err := alert.CreateIfNotExists(c.ID, model.MetricTypeTemperature, model.AlertSeverityWarning, "High temp", "39.6°C exceeds warning threshold")

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND status = 'active'", c.ID).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestCreateIfNotExists_NoDuplicate(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Cow2", model.CowStatusInFarm)

		err1 := alert.CreateIfNotExists(c.ID, model.MetricTypeTemperature, model.AlertSeverityWarning, "High temp", "39.6°C")
		err2 := alert.CreateIfNotExists(c.ID, model.MetricTypeTemperature, model.AlertSeverityCritical, "High temp", "40.6°C")

		require.NoError(t, err1)
		require.NoError(t, err2)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND metric_key = ? AND status = 'active'", c.ID, model.MetricTypeTemperature).Count(&count)
		assert.Equal(t, int64(1), count)
	})
}

func TestResolveIfExists_ResolvesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Cow3", model.CowStatusInFarm)
		testhelper.SeedAlertForMetric(t, tx, c.ID, model.MetricTypeTemperature, model.AlertSeverityWarning, model.AlertStatusActive)

		err := alert.ResolveIfExists(c.ID, model.MetricTypeTemperature)

		require.NoError(t, err)
		var count int64
		tx.Model(&model.Alert{}).Where("cow_id = ? AND status = 'active'", c.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestResolveIfExists_NoopWhenNoAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Cow4", model.CowStatusInFarm)

		err := alert.ResolveIfExists(c.ID, model.MetricTypeTemperature)

		require.NoError(t, err)
	})
}
