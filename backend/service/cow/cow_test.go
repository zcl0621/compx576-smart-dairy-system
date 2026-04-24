package cow_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cowdto "github.com/zcl0621/compx576-smart-dairy-system/dto/cow"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/cow"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func nowFunc() time.Time { return time.Now() }

func TestCowCreate_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		err := cow.CowCreateService(&cowdto.CreateRequest{
			Name:   "Bessie",
			Tag:    "T-001",
			Age:    3,
			Status: model.CowStatusInFarm,
		})

		require.NoError(t, err)
	})
}

func TestCowList_FilterByStatus(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedCow(t, tx, "InFarm1", model.CowStatusInFarm)
		testhelper.SeedCow(t, tx, "Sold1", model.CowStatusSold)

		resp, err := cow.CowListService(&cowdto.ListQuery{Status: string(model.CowStatusInFarm)})

		require.NoError(t, err)
		for _, item := range resp.List {
			assert.Equal(t, model.CowStatusInFarm, item.Status)
		}
	})
}

func TestCowList_FilterByCondition(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c1 := testhelper.SeedCow(t, tx, "Critical1", model.CowStatusInFarm)
		testhelper.SeedCow(t, tx, "Normal1", model.CowStatusInFarm)
		testhelper.SeedAlertForMetric(t, tx, c1.ID, model.MetricTypeTemperature, model.AlertSeverityCritical, model.AlertStatusActive)

		resp, err := cow.CowListService(&cowdto.ListQuery{Condition: string(model.CowConditionCritical)})

		require.NoError(t, err)
		require.GreaterOrEqual(t, int(resp.Total), 1)
		for _, item := range resp.List {
			assert.Equal(t, model.CowConditionCritical, item.Condition)
		}
	})
}

func TestCowList_SortByCondition(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		offline := testhelper.SeedCow(t, tx, "OfflineCow", model.CowStatusInFarm)
		critical := testhelper.SeedCow(t, tx, "CriticalCow", model.CowStatusInFarm)
		testhelper.SeedCow(t, tx, "NormalCow", model.CowStatusInFarm)
		testhelper.SeedAlertForMetric(t, tx, offline.ID, model.MetricTypeDevice, model.AlertSeverityOffline, model.AlertStatusActive)
		testhelper.SeedAlertForMetric(t, tx, critical.ID, model.MetricTypeTemperature, model.AlertSeverityCritical, model.AlertStatusActive)

		resp, err := cow.CowListService(&cowdto.ListQuery{Sort: "condition"})

		require.NoError(t, err)
		require.GreaterOrEqual(t, len(resp.List), 3)
		conditionOrder := map[model.CowCondition]int{
			model.CowConditionCritical: 1,
			model.CowConditionWarning:  2,
			model.CowConditionNormal:   3,
			model.CowConditionOffline:  4,
		}
		for i := 1; i < len(resp.List); i++ {
			prev := conditionOrder[resp.List[i-1].Condition]
			curr := conditionOrder[resp.List[i].Condition]
			assert.LessOrEqual(t, prev, curr, "condition order violated at index %d", i)
		}
	})
}

func TestCowInfo_WithLatestWeight(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "WeightCow", model.CowStatusInFarm)
		testhelper.SeedMetric(t, tx, c.ID, model.MetricTypeWeight, 450.5, nowFunc())

		resp, err := cow.CowInfoService(&cowdto.InfoQuery{ID: c.ID})

		require.NoError(t, err)
		require.NotNil(t, resp.Weight)
		assert.InDelta(t, 450.5, *resp.Weight, 0.01)
	})
}

func TestCowInfo_NoMetrics(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "NoMetricCow", model.CowStatusInFarm)

		resp, err := cow.CowInfoService(&cowdto.InfoQuery{ID: c.ID})

		require.NoError(t, err)
		assert.Nil(t, resp.Weight)
	})
}

func TestCowUpdate_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "OldName", model.CowStatusInFarm)

		err := cow.CowUpdateService(&cowdto.UpdateRequest{
			ID:     c.ID,
			Name:   "NewName",
			Tag:    c.Tag,
			Age:    5,
			Status: model.CowStatusInFarm,
		})
		require.NoError(t, err)

		resp, err := cow.CowInfoService(&cowdto.InfoQuery{ID: c.ID})
		require.NoError(t, err)
		assert.Equal(t, "NewName", resp.Name)
		assert.Equal(t, 5, resp.Age)
	})
}

func TestCowList_SearchByName(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedCow(t, tx, "Findme", model.CowStatusInFarm)
		testhelper.SeedCow(t, tx, "Other", model.CowStatusInFarm)

		resp, err := cow.CowListService(&cowdto.ListQuery{Name: "Findme"})

		require.NoError(t, err)
		require.GreaterOrEqual(t, int(resp.Total), 1)
		for _, item := range resp.List {
			assert.Contains(t, item.Name, "Findme")
		}
	})
}

func TestCowList_SortByUpdatedAt(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		testhelper.SeedCow(t, tx, "SortCow1", model.CowStatusInFarm)
		testhelper.SeedCow(t, tx, "SortCow2", model.CowStatusInFarm)

		// updated_at is default, just check no error
		resp, err := cow.CowListService(&cowdto.ListQuery{Sort: "updated_at"})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, int(resp.Total), 2)
		// check desc order
		for i := 1; i < len(resp.List); i++ {
			assert.False(t, resp.List[i].UpdatedAt.After(resp.List[i-1].UpdatedAt),
				"updated_at order violated at index %d", i)
		}
	})
}

func TestCowInfo_NotFound(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		_, err := cow.CowInfoService(&cowdto.InfoQuery{ID: "nonexistent-cow-id"})

		assert.Error(t, err)
	})
}

func TestCowCreate_DuplicateTag(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "First", model.CowStatusInFarm)

		err := cow.CowCreateService(&cowdto.CreateRequest{
			Name:   "Second",
			Tag:    c.Tag, // keep same tag
			Age:    2,
			Status: model.CowStatusInFarm,
		})

		assert.Error(t, err)
	})
}

func TestCowUpdate_NotFound(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		err := cow.CowUpdateService(&cowdto.UpdateRequest{
			ID:     "nonexistent-cow-id",
			Name:   "Ghost",
			Tag:    "GHOST-TAG",
			Status: model.CowStatusInFarm,
		})

		assert.Error(t, err)
	})
}
