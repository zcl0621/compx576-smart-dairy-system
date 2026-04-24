package cow

import (
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	cowdto "github.com/zcl0621/compx576-smart-dairy-system/dto/cow"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"gorm.io/gorm"
)

func CowListService(r *cowdto.ListQuery) (*cowdto.ListResponse, error) {
	db := cowBaseQuery()

	switch r.Sort {
	case "condition":
		db = model.OrderByCowCondition(db)
	default:
		db = model.OrderByUpdatedAt(db)
	}

	if r.Name != "" {
		db = db.Where("cows.name LIKE ?", "%"+r.Name+"%")
	}
	if r.Status != "" {
		db = db.Where("cows.status = ?", r.Status)
	}
	if r.Condition != "" {
		db = applyConditionFilter(db, model.CowCondition(r.Condition))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []cowdto.ListItem
	if err := db.Offset(r.GetOffset()).Limit(r.GetLimit()).Find(&list).Error; err != nil {
		return nil, err
	}

	return &cowdto.ListResponse{
		List: list,
		PageResponse: common.PageResponse{
			Page:       r.GetPage(),
			Total:      total,
			TotalPages: r.GetTotalPages(total),
		},
	}, nil
}

func CowInfoService(r *cowdto.InfoQuery) (*cowdto.InfoResponse, error) {
	var response cowdto.InfoResponse
	err := cowBaseQuery().
		Where("cows.id = ?", r.ID).
		First(&response).Error
	if err != nil {
		return nil, err
	}

	latestOf := func(metricType model.MetricType) *float64 {
		var row struct{ MetricValue float64 }
		if err := pg.DB.Model(&model.Metric{}).
			Select("metric_value").
			Where("cow_id = ? AND metric_type = ?", r.ID, metricType).
			Order("created_at desc").
			Take(&row).Error; err == nil {
			v := row.MetricValue
			return &v
		}
		return nil
	}

	response.Weight = latestOf(model.MetricTypeWeight)
	response.Temperature = latestOf(model.MetricTypeTemperature)
	response.HeartRate = latestOf(model.MetricTypeHeartRate)
	response.BloodOxygen = latestOf(model.MetricTypeBloodOxygen)
	response.MilkAmount = latestOf(model.MetricTypeMilkAmount)

	return &response, nil
}

func CowCreateService(r *cowdto.CreateRequest) error {
	item := model.Cow{
		Name:       r.Name,
		Tag:        r.Tag,
		Age:        r.Age,
		CanMilking: r.CanMilking,
		Status:     r.Status,
	}
	return pg.DB.Create(&item).Error
}

func CowUpdateService(r *cowdto.UpdateRequest) error {
	var item model.Cow
	if err := pg.DB.Where("id = ?", r.ID).First(&item).Error; err != nil {
		return err
	}
	item.Name = r.Name
	item.Tag = r.Tag
	item.Age = r.Age
	item.CanMilking = r.CanMilking
	item.Status = r.Status
	return pg.DB.Save(&item).Error
}

func cowBaseQuery() *gorm.DB {
	return pg.DB.Model(&model.Cow{}).Select("cows.*", model.ConditionSubQuery()+" AS condition")
}

func applyConditionFilter(db *gorm.DB, condition model.CowCondition) *gorm.DB {
	if condition == model.CowConditionNormal {
		return db.Where("NOT EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = 'active' AND alerts.deleted_at IS NULL)")
	}
	return db.Where("EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = 'active' AND alerts.severity = ? AND alerts.deleted_at IS NULL)", condition)
}
