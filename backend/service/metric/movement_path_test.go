package metric_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/metric"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestMovementPathService_ReturnsPathPoints(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "PathCow", model.CowStatusInFarm, model.CowConditionNormal)
		now := time.Now()

		// seed 3 non-collinear positions: A -> B -> C
		// values chosen to survive numeric(10,2) rounding and differ by > epsilon after DP
		testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLatitude, -37.78, now.Add(-3*time.Hour))
		testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLongitude, 175.27, now.Add(-3*time.Hour))
		testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLatitude, -37.79, now.Add(-2*time.Hour))
		testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLongitude, 175.30, now.Add(-2*time.Hour))
		testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLatitude, -37.80, now.Add(-1*time.Hour))
		testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLongitude, 175.29, now.Add(-1*time.Hour))

		resp, err := metric.MovementPathService(&metric.MetricQuery{
			CowID:       cow.ID,
			MetricRange: model.MetricRange24H,
		})

		require.NoError(t, err)
		assert.Equal(t, cow.ID, resp.CowID)
		assert.GreaterOrEqual(t, len(resp.Points), 3)
		assert.InDelta(t, -37.78, resp.Points[0].Lat, 0.001)
	})
}

func TestMovementPathService_DetectsStay(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "StayCow", model.CowStatusInFarm, model.CowConditionNormal)
		now := time.Now()

		// seed 3 points at same rounded location (within 5m) over 10 min
		// use values that round to same numeric(10,2) value
		for i := 0; i < 3; i++ {
			offset := time.Duration(i*5) * time.Minute
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLatitude, -37.78, now.Add(-30*time.Minute+offset))
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLongitude, 175.27, now.Add(-30*time.Minute+offset))
		}
		// seed 1 point far away (> 5m after rounding)
		testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLatitude, -37.90, now.Add(-10*time.Minute))
		testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeLongitude, 175.40, now.Add(-10*time.Minute))

		resp, err := metric.MovementPathService(&metric.MetricQuery{
			CowID:       cow.ID,
			MetricRange: model.MetricRange24H,
		})

		require.NoError(t, err)
		// 3 close points merge into 1 stay + 1 distant = 2 points
		assert.Equal(t, 2, len(resp.Points))
		// first point should have stay_seconds > 0
		assert.Greater(t, resp.Points[0].StaySeconds, int64(0))
	})
}

func TestMovementPathService_EmptyData(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "EmptyCow", model.CowStatusInFarm, model.CowConditionNormal)

		resp, err := metric.MovementPathService(&metric.MetricQuery{
			CowID:       cow.ID,
			MetricRange: model.MetricRange24H,
		})

		require.NoError(t, err)
		assert.Equal(t, 0, len(resp.Points))
	})
}
