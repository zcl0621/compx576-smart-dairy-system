package metric

import (
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	metricdto "github.com/zcl0621/compx576-smart-dairy-system/dto/metric"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"gorm.io/gorm"
)

func MetricListService(r *metricdto.ListQuery) (*metricdto.ListResponse, error) {
	db := metricListBaseQuery()
	db = applyMetricListFilters(db, r.CowID, r.MetricType)
	db = db.Order("metrics.created_at desc")

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []metricdto.ListItem
	if err := db.Offset(r.GetOffset()).Limit(r.GetLimit()).Find(&list).Error; err != nil {
		return nil, err
	}

	response := metricdto.ListResponse{
		List: list,
		PageResponse: common.PageResponse{
			Page:       r.GetPage(),
			Total:      total,
			TotalPages: r.GetTotalPages(total),
		},
	}

	return &response, nil
}

func metricListBaseQuery() *gorm.DB {
	return pg.DB.Model(&model.Metric{}).
		Joins("LEFT JOIN cows ON cows.id = metrics.cow_id AND cows.deleted_at IS NULL").
		Select(
			"metrics.id",
			"metrics.cow_id",
			"cows.name as cow_name",
			"metrics.source",
			"metrics.metric_type",
			"metrics.metric_value",
			"metrics.unit",
			"metrics.created_at",
		)
}

func applyMetricListFilters(db *gorm.DB, cowID string, metricType string) *gorm.DB {
	if cowID != "" {
		db = db.Where("metrics.cow_id = ?", cowID)
	}
	if metricType != "" {
		db = db.Where("metrics.metric_type = ?", metricType)
	}
	return db
}
