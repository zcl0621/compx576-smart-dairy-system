package alert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	alertdto "github.com/zcl0621/compx576-smart-dairy-system/dto/alert"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/alert"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestAlertList_All(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "ListCow", model.CowStatusInFarm)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityCritical, model.AlertStatusActive)

		resp, err := alert.ListService(&alertdto.ListQuery{})

		require.NoError(t, err)
		assert.Equal(t, int64(2), int64(resp.Total))
	})
}

func TestAlertList_FilterByCow(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c1 := testhelper.SeedCow(t, tx, "Cow1", model.CowStatusInFarm)
		c2 := testhelper.SeedCow(t, tx, "Cow2", model.CowStatusInFarm)
		testhelper.SeedAlert(t, tx, c1.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		testhelper.SeedAlert(t, tx, c2.ID, model.AlertSeverityCritical, model.AlertStatusActive)

		resp, err := alert.ListService(&alertdto.ListQuery{CowID: c1.ID})

		require.NoError(t, err)
		require.Equal(t, int64(1), resp.Total)
		assert.Equal(t, c1.ID, resp.List[0].CowID)
	})
}

func TestAlertList_FilterBySeverity(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "SevCow", model.CowStatusInFarm)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityCritical, model.AlertStatusActive)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityOffline, model.AlertStatusActive)

		resp, err := alert.ListService(&alertdto.ListQuery{Severity: string(model.AlertSeverityCritical)})

		require.NoError(t, err)
		require.Equal(t, int64(1), resp.Total)
		assert.Equal(t, model.AlertSeverityCritical, resp.List[0].Severity)
	})
}

func TestAlertList_ExcludesResolved(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "ResolvedCow", model.CowStatusInFarm)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusResolved)

		resp, err := alert.ListService(&alertdto.ListQuery{CowID: c.ID})

		require.NoError(t, err)
		assert.Equal(t, int64(0), resp.Total)
	})
}

func TestAlertList_Pagination(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "PagCow", model.CowStatusInFarm)
		for i := 0; i < 5; i++ {
			testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		}

		q := &alertdto.ListQuery{CowID: c.ID}
		q.PageQuery = common.PageQuery{Page: 1, PageSize: 2}

		resp, err := alert.ListService(q)

		require.NoError(t, err)
		assert.Equal(t, int64(5), resp.Total)
		assert.Len(t, resp.List, 2)
		assert.Equal(t, 3, resp.TotalPages)
	})
}

func TestAlertSummary_Counts(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "SummaryCow", model.CowStatusInFarm)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityWarning, model.AlertStatusActive)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityCritical, model.AlertStatusActive)
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityOffline, model.AlertStatusActive)
		// don't count resolved alerts
		testhelper.SeedAlert(t, tx, c.ID, model.AlertSeverityCritical, model.AlertStatusResolved)

		resp, err := alert.SummaryService()

		require.NoError(t, err)
		// totals are exact now, db is clean
		assert.Equal(t, int64(4), resp.Active)
		assert.Equal(t, int64(2), resp.Warning)
		assert.Equal(t, int64(1), resp.Critical)
		assert.Equal(t, int64(1), resp.Offline)
	})
}
