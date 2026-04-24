package metric_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metricdto "github.com/zcl0621/compx576-smart-dairy-system/dto/metric"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/metric"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestMetricListService_ReturnsAll(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "AllCow", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.5, now)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeHeartRate, 72.0, now)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeBloodOxygen, 98.0, now)

		resp, err := metric.MetricListService(&metricdto.ListQuery{})

		require.NoError(t, err)
		assert.Equal(t, int64(3), resp.Total)
	})
}

func TestMetricListService_FilterByCowID(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c1 := testhelper.SeedCow(t, tx, "Cow1", model.CowStatusInFarm)
		c2 := testhelper.SeedCow(t, tx, "Cow2", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c1.ID, model.MetricTypeTemperature, 38.5, now)
		testhelper.SeedMetric(t, tx, c1.ID, model.MetricTypeHeartRate, 72.0, now)
		testhelper.SeedMetric(t, tx, c2.ID, model.MetricTypeTemperature, 39.0, now)

		resp, err := metric.MetricListService(&metricdto.ListQuery{CowID: c1.ID})

		require.NoError(t, err)
		require.Equal(t, int64(2), resp.Total)
		for _, item := range resp.List {
			assert.Equal(t, c1.ID, item.CowID)
		}
	})
}

func TestMetricListService_FilterByMetricType(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "TypeCow", model.CowStatusInFarm)
		now := time.Now()
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.5, now)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, 38.8, now)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeHeartRate, 72.0, now)

		resp, err := metric.MetricListService(&metricdto.ListQuery{MetricType: string(model.MetricTypeTemperature)})

		require.NoError(t, err)
		require.Equal(t, int64(2), resp.Total)
		for _, item := range resp.List {
			assert.Equal(t, model.MetricTypeTemperature, item.MetricType)
		}
	})
}

func TestMetricListService_Pagination(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "PagCow", model.CowStatusInFarm)
		now := time.Now()
		for i := 0; i < 5; i++ {
			testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeTemperature, float64(38+i), now)
		}

		q := &metricdto.ListQuery{CowID: c.ID}
		q.PageQuery = common.PageQuery{Page: 1, PageSize: 2}

		resp, err := metric.MetricListService(q)

		require.NoError(t, err)
		assert.Equal(t, int64(5), resp.Total)
		assert.Len(t, resp.List, 2)
		assert.Equal(t, 3, resp.TotalPages)
	})
}
