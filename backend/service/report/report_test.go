package report_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	reportdto "github.com/zcl0621/compx576-smart-dairy-system/dto/report"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/report"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestReportList_ReturnsSeededReports(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Bessie", model.CowStatusInFarm, model.CowConditionNormal)
		testhelper.SeedReport(t, tx, c.ID)
		testhelper.SeedReport(t, tx, c.ID)

		q := &reportdto.ListQuery{}
		q.PageQuery = common.PageQuery{Page: 1, PageSize: 20}

		resp, err := report.ListService(q)

		require.NoError(t, err)
		assert.Equal(t, int64(2), resp.Total)
		assert.Len(t, resp.List, 2)
	})
}

func TestReportList_CowNameJoin(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Daisy", model.CowStatusInFarm, model.CowConditionNormal)
		testhelper.SeedReport(t, tx, c.ID)

		q := &reportdto.ListQuery{}
		q.PageQuery = common.PageQuery{Page: 1, PageSize: 20}

		resp, err := report.ListService(q)

		require.NoError(t, err)
		// find seeded report and check cow_name
		var found *reportdto.ReportItem
		for i := range resp.List {
			if resp.List[i].CowID == c.ID {
				found = &resp.List[i]
				break
			}
		}
		require.NotNil(t, found, "seeded report not found in list")
		assert.Equal(t, "Daisy", found.CowName)
	})
}

func TestReportList_OrderedByCreatedAtDesc(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "Molly", model.CowStatusInFarm, model.CowConditionNormal)
		r1 := testhelper.SeedReport(t, tx, c.ID)
		r2 := testhelper.SeedReport(t, tx, c.ID)

		q := &reportdto.ListQuery{}
		q.PageQuery = common.PageQuery{Page: 1, PageSize: 20}

		resp, err := report.ListService(q)

		require.NoError(t, err)

		// get r1 and r2 pos
		pos := map[string]int{}
		for i, item := range resp.List {
			if item.ID == r1.ID || item.ID == r2.ID {
				pos[item.ID] = i
			}
		}
		require.Len(t, pos, 2, "both seeded reports should appear in result")
		// newer r2 comes before older r1
		assert.Less(t, pos[r2.ID], pos[r1.ID], "newer report should appear first")
	})
}

func TestReportList_Pagination(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "PagCow", model.CowStatusInFarm, model.CowConditionNormal)
		for i := 0; i < 3; i++ {
			testhelper.SeedReport(t, tx, c.ID)
		}

		// page 1 size 2
		q := &reportdto.ListQuery{}
		q.PageQuery = common.PageQuery{Page: 1, PageSize: 2}

		resp, err := report.ListService(q)

		require.NoError(t, err)
		assert.LessOrEqual(t, len(resp.List), 2)
		assert.Equal(t, resp.Page, 1)
	})
}

func TestReportLatest_Success(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		c := testhelper.SeedCow(t, tx, "LatestCow", model.CowStatusInFarm, model.CowConditionNormal)
		r1 := testhelper.SeedReport(t, tx, c.ID)
		r2 := testhelper.SeedReport(t, tx, c.ID)

		item, err := report.LatestService(&reportdto.LatestQuery{CowID: c.ID})

		require.NoError(t, err)
		// get latest report
		assert.Equal(t, r2.ID, item.ID)
		assert.Equal(t, c.ID, item.CowID)
		_ = r1
	})
}

func TestReportLatest_NotFound(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		item, err := report.LatestService(&reportdto.LatestQuery{CowID: "nonexistent-id"})

		assert.Error(t, err)
		assert.Nil(t, item)
	})
}
