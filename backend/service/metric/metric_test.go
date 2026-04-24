package metric_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cowdto "github.com/zcl0621/compx576-smart-dairy-system/dto/cow"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/metric"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

// --- temp ---

func TestTemperatureService_ReturnsData(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "TempCow", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.5, time.Now())

		resp, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, c.ID, resp.CowID)
		assert.Len(t, resp.Series, 1)
		require.NotNil(t, resp.Summary.Current)
		assert.InDelta(t, 38.5, *resp.Summary.Current, 0.01)
		assert.Equal(t, model.ReportMetricStatusNormal, resp.Summary.Status)
	})
}

func TestTemperatureService_StatusWarning(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "TempWarn", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 39.6, time.Now())

		resp, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusWarning, resp.Summary.Status)
	})
}

func TestTemperatureService_StatusCritical(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "TempCrit", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 40.6, time.Now())

		resp, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusCritical, resp.Summary.Status)
	})
}

func TestTemperatureService_StatusOffline(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "TempOffline", model.CowStatusInFarm)
		// no metric here

		resp, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusOffline, resp.Summary.Status)
		assert.Empty(t, resp.Series)
		assert.Nil(t, resp.Summary.Current)
	})
}

func TestTemperatureService_BadRange(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "RangeCow", model.CowStatusInFarm)

		_, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "invalid"})

		assert.ErrorIs(t, err, metric.ErrBadMetricRange)
	})
}

// --- heart ---

func TestHeartRateService_StatusNormal(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HRNormal", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeHeartRate, 70.0, time.Now())

		resp, err := metric.HeartRateService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusNormal, resp.Summary.Status)
	})
}

func TestHeartRateService_StatusWarning_High(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HRWarnHigh", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeHeartRate, 95.0, time.Now())

		resp, err := metric.HeartRateService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusWarning, resp.Summary.Status)
	})
}

func TestHeartRateService_StatusCritical_High(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HRCritHigh", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeHeartRate, 105.0, time.Now())

		resp, err := metric.HeartRateService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusCritical, resp.Summary.Status)
	})
}

func TestHeartRateService_StatusCritical_Low(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HRCritLow", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeHeartRate, 37.0, time.Now())

		resp, err := metric.HeartRateService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusCritical, resp.Summary.Status)
	})
}

// --- oxygen ---

func TestBloodOxygenService_StatusNormal(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "BONormal", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeBloodOxygen, 97.0, time.Now())

		resp, err := metric.BloodOxygenService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusNormal, resp.Summary.Status)
	})
}

func TestBloodOxygenService_StatusWarning(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "BOWarning", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeBloodOxygen, 93.0, time.Now())

		resp, err := metric.BloodOxygenService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusWarning, resp.Summary.Status)
	})
}

func TestBloodOxygenService_StatusCritical(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "BOCritical", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeBloodOxygen, 87.0, time.Now())

		resp, err := metric.BloodOxygenService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusCritical, resp.Summary.Status)
	})
}

// --- milk ---

func TestMilkAmountService_SumAndAvg(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MilkCow", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeMilkAmount, 10.0, now.Add(-2*time.Hour))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeMilkAmount, 12.0, now.Add(-1*time.Hour))

		resp, err := metric.MilkAmountService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, int64(2), resp.Summary.SessionCount)
		assert.InDelta(t, 22.0, resp.Summary.Total, 0.01)
		assert.InDelta(t, 11.0, resp.Summary.AvgPerSession, 0.01)
	})
}

// --- move ---

func TestMovementService_PairLatLng(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MoveCow", model.CowStatusInFarm)
		now := time.Now()
		// pair 2 lat lng points, 1 min apart
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeLatitude, -37.78, now.Add(-10*time.Minute))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeLongitude, 175.28, now.Add(-10*time.Minute))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeLatitude, -37.79, now.Add(-5*time.Minute))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeLongitude, 175.29, now.Add(-5*time.Minute))

		resp, err := metric.MovementService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, int64(2), resp.Summary.PointCount)
		assert.NotEmpty(t, resp.Series)
		assert.Greater(t, resp.Summary.DistanceM, 0.0)
		assert.Equal(t, model.ReportMetricStatusNormal, resp.Summary.Status)
	})
}

func TestMovementService_SinglePoint_Warning(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MoveSingle", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeLatitude, -37.78, now)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeLongitude, 175.28, now)

		resp, err := metric.MovementService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusWarning, resp.Summary.Status)
	})
}

func TestMovementService_NoData_Offline(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MoveNone", model.CowStatusInFarm)

		resp, err := metric.MovementService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusOffline, resp.Summary.Status)
		assert.Empty(t, resp.Series)
	})
}

// --- ranges ---

func TestTemperatureService_AllRange(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "RangeAll", model.CowStatusInFarm)
		// seed old metric so all gets it
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.0, time.Now().Add(-48*time.Hour))

		resp24h, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})
		require.NoError(t, err)

		respAll, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "all"})
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(respAll.Series), len(resp24h.Series))
	})
}

func TestTemperatureService_7DRange(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Range7D", model.CowStatusInFarm)
		// inside 7d
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.2, time.Now().Add(-3*24*time.Hour))
		// outside 7d, should show in all only
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.5, time.Now().Add(-10*24*time.Hour))

		resp7d, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "7d"})
		require.NoError(t, err)

		respAll, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "all"})
		require.NoError(t, err)

		assert.Equal(t, 1, len(resp7d.Series))
		assert.Equal(t, 2, len(respAll.Series))
	})
}

func TestTemperatureService_30DRange(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Range30D", model.CowStatusInFarm)
		// inside 30d
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.1, time.Now().Add(-15*24*time.Hour))
		// outside 30d
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.8, time.Now().Add(-40*24*time.Hour))

		resp30d, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "30d"})
		require.NoError(t, err)

		respAll, err := metric.TemperatureService(&cowdto.MetricQuery{CowID: c.ID, Range: "all"})
		require.NoError(t, err)

		assert.Equal(t, 1, len(resp30d.Series))
		assert.Equal(t, 2, len(respAll.Series))
	})
}

func TestHeartRateService_StatusWarning_Low(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "HRWarnLow", model.CowStatusInFarm)
		// 38-48 is low warning range (< HRWarnLow=48 but >= HRCritLow=38)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeHeartRate, 47.0, time.Now())

		resp, err := metric.HeartRateService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, model.ReportMetricStatusWarning, resp.Summary.Status)
	})
}

func TestMilkAmountService_NoData(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "DryMilkCow", model.CowStatusInFarm)

		resp, err := metric.MilkAmountService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, int64(0), resp.Summary.SessionCount)
		assert.InDelta(t, 0.0, resp.Summary.Total, 0.001)
		assert.InDelta(t, 0.0, resp.Summary.AvgPerSession, 0.001)
		assert.Nil(t, resp.UpdatedAt)
	})
}

func TestMovementService_PairBeyondTolerance_Dropped(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MoveTolerCow", model.CowStatusInFarm)
		now := time.Now()
		// lat at t=0, lng comes 10 min later, drop it
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeLatitude, -37.78, now.Add(-20*time.Minute))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeLongitude, 175.28, now.Add(-8*time.Minute))

		resp, err := metric.MovementService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		// no valid pair, goes offline
		assert.Equal(t, model.ReportMetricStatusOffline, resp.Summary.Status)
		assert.Equal(t, int64(0), resp.Summary.PointCount)
	})
}

// --- weight ---

func TestWeightService_ReturnsData(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "WeightCow", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeWeight, 520.0, time.Now())

		resp, err := metric.WeightService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Equal(t, c.ID, resp.CowID)
		assert.Len(t, resp.Series, 1)
		require.NotNil(t, resp.Summary.Current)
		assert.InDelta(t, 520.0, *resp.Summary.Current, 0.01)
		require.NotNil(t, resp.UpdatedAt)
	})
}

func TestWeightService_NoData(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "EmptyWeightCow", model.CowStatusInFarm)

		resp, err := metric.WeightService(&cowdto.MetricQuery{CowID: c.ID, Range: "24h"})

		require.NoError(t, err)
		assert.Empty(t, resp.Series)
		assert.Nil(t, resp.Summary.Current)
		assert.Nil(t, resp.Summary.Avg)
		assert.Nil(t, resp.UpdatedAt)
	})
}

func TestWeightService_StatsCorrect(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "MultiWeightCow", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeWeight, 500.0, now.Add(-2*time.Hour))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeWeight, 520.0, now.Add(-1*time.Hour))
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeWeight, 510.0, now)

		resp, err := metric.WeightService(&cowdto.MetricQuery{CowID: c.ID, Range: "all"})

		require.NoError(t, err)
		assert.Len(t, resp.Series, 3)
		require.NotNil(t, resp.Summary.Current)
		assert.InDelta(t, 510.0, *resp.Summary.Current, 0.01)
		require.NotNil(t, resp.Summary.Avg)
		assert.InDelta(t, 510.0, *resp.Summary.Avg, 0.01) // avg of 500 520 510
		require.NotNil(t, resp.Summary.Min)
		assert.InDelta(t, 500.0, *resp.Summary.Min, 0.01)
		require.NotNil(t, resp.Summary.Max)
		assert.InDelta(t, 520.0, *resp.Summary.Max, 0.01)
	})
}
