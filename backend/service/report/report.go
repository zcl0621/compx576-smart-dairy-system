package report

import (
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	reportdto "github.com/zcl0621/compx576-smart-dairy-system/dto/report"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"gorm.io/gorm"
)

func ListService(r *reportdto.ListQuery) (*reportdto.ListResponse, error) {
	db := reportBaseQuery()
	db = db.Order("reports.created_at desc")

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []reportdto.ReportItem
	if err := db.Offset(r.GetOffset()).Limit(r.GetLimit()).Find(&list).Error; err != nil {
		return nil, err
	}

	response := reportdto.ListResponse{
		List: list,
		PageResponse: common.PageResponse{
			Page:       r.GetPage(),
			Total:      total,
			TotalPages: r.GetTotalPages(total),
		},
	}

	return &response, nil
}

func LatestService(r *reportdto.LatestQuery) (*reportdto.ReportItem, error) {
	db := reportBaseQuery()
	db = db.Where("reports.cow_id = ?", r.CowID)
	db = db.Order("reports.created_at desc")

	var item reportdto.ReportItem
	if err := db.First(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func reportBaseQuery() *gorm.DB {
	return pg.DB.Model(&model.Report{}).
		Select(
			"reports.id",
			"reports.created_at",
			"reports.updated_at",
			"reports.cow_id",
			"(SELECT name FROM cows WHERE cows.id = reports.cow_id) AS cow_name",
			"reports.period_start",
			"reports.period_end",
			"reports.summary",
			"reports.score",
			"reports.details_json AS details",
		)
}
