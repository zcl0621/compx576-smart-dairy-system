package dashboard_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	dashboarddto "github.com/zcl0621/compx576-smart-dairy-system/dto/dashboard"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/dashboard"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func nowFunc() time.Time { return time.Now() }

func TestDashboardSummary_Counts(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		n1 := testhelper.SeedCow(t, tx, "N1", model.CowStatusInFarm)
		n2 := testhelper.SeedCow(t, tx, "N2", model.CowStatusInFarm)
		w1 := testhelper.SeedCow(t, tx, "W1", model.CowStatusInFarm)
		c1 := testhelper.SeedCow(t, tx, "C1", model.CowStatusInFarm)
		o1 := testhelper.SeedCow(t, tx, "O1", model.CowStatusInFarm)
		// sold cow must not show in count
		testhelper.SeedCow(t, tx, "Sold", model.CowStatusSold)

		// suppress "declared but not used" for n1, n2 (they stay normal — no alerts needed)
		_ = n1
		_ = n2

		// set conditions via active alerts
		testhelper.SeedAlertForMetric(t, tx, w1.ID, model.MetricTypeHeartRate, model.AlertSeverityWarning, model.AlertStatusActive)
		testhelper.SeedAlertForMetric(t, tx, c1.ID, model.MetricTypeTemperature, model.AlertSeverityCritical, model.AlertStatusActive)
		testhelper.SeedAlertForMetric(t, tx, o1.ID, model.MetricTypeDevice, model.AlertSeverityOffline, model.AlertStatusActive)

		resp, err := dashboard.SummaryService()

		require.NoError(t, err)
		assert.Equal(t, int64(5), resp.TotalCows)
		assert.Equal(t, int64(2), resp.Normal)
		assert.Equal(t, int64(1), resp.Warning)
		assert.Equal(t, int64(1), resp.Critical)
		assert.Equal(t, int64(1), resp.Offline)
	})
}

func TestDashboardSummary_SoldExcluded(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		// get base count with no sold cows
		baseline, err := dashboard.SummaryService()
		require.NoError(t, err)

		testhelper.SeedCow(t, tx, "SoldExtra", model.CowStatusSold)

		resp, err := dashboard.SummaryService()
		require.NoError(t, err)

		// total must not go up after sold cow
		assert.Equal(t, baseline.TotalCows, resp.TotalCows)
	})
}

func TestDashboardList_WithMetrics(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MetricCow", model.CowStatusInFarm)
		now := nowFunc()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.5, now)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeHeartRate, 72.0, now)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeBloodOxygen, 98.0, now)

		resp, err := dashboard.ListService(&dashboarddto.ListQuery{})

		require.NoError(t, err)
		// find seeded cow
		var found *dashboarddto.ListItem
		for i := range resp.List {
			if resp.List[i].ID == c.ID {
				found = &resp.List[i]
				break
			}
		}
		require.NotNil(t, found, "seeded cow not found in dashboard list")
		require.NotNil(t, found.Temperature)
		require.NotNil(t, found.HeartRate)
		require.NotNil(t, found.BloodOxygen)
		assert.InDelta(t, 38.5, *found.Temperature, 0.01)
		assert.InDelta(t, 72.0, *found.HeartRate, 0.01)
		assert.InDelta(t, 98.0, *found.BloodOxygen, 0.01)
	})
}

func TestDashboardList_NoMetrics(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "NoMetricCow", model.CowStatusInFarm)

		resp, err := dashboard.ListService(&dashboarddto.ListQuery{})

		require.NoError(t, err)
		var found *dashboarddto.ListItem
		for i := range resp.List {
			if resp.List[i].ID == c.ID {
				found = &resp.List[i]
				break
			}
		}
		require.NotNil(t, found)
		assert.Nil(t, found.Temperature)
		assert.Nil(t, found.HeartRate)
		assert.Nil(t, found.BloodOxygen)
	})
}

func TestDashboardList_AlertMessage(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "AlertCow", model.CowStatusInFarm)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		a := testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityCritical, model.AlertStatusActive)

		resp, err := dashboard.ListService(&dashboarddto.ListQuery{})

		require.NoError(t, err)
		var found *dashboarddto.ListItem
		for i := range resp.List {
			if resp.List[i].ID == c.ID {
				found = &resp.List[i]
				break
			}
		}
		require.NotNil(t, found)
		require.NotNil(t, found.AlertMessage)
		// critical msg comes first
		assert.Equal(t, a.Message, *found.AlertMessage)
	})
}
