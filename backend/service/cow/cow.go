package cow

import (
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/cow"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

func CowListService(r *cow.ListQuery) (*cow.ListResponse, error) {
	db := pg.DB.Model(&model.Cow{})
	switch r.Sort {
	case "updated_at":
		db = model.OrderByUpdatedAt(db)
	case "condition":
		db = model.OrderByCowCondition(db)
	}
	if r.Name != "" {
		db = db.Where("name LIKE ?", "%"+r.Name+"%")
	}
	if r.Status != "" {
		db = db.Where(&model.Cow{Status: model.CowStatus(r.Status)})
	}
	if r.Condition != "" {
		db = db.Where(&model.Cow{Condition: model.CowCondition(r.Condition)})
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []cow.ListItem
	if err := db.Offset(r.GetOffset()).Limit(r.GetLimit()).Find(&list).Error; err != nil {
		return nil, err
	}

	response := cow.ListResponse{
		List: list,
		PageResponse: common.PageResponse{
			Page:       r.GetPage(),
			Total:      total,
			TotalPages: r.GetTotalPages(total),
		},
	}

	return &response, nil
}

func CowInfoService(r *cow.InfoQuery) (*cow.InfoResponse, error) {
	var response cow.InfoResponse
	err := pg.DB.Model(&model.Cow{}).
		Where("id = ?", r.ID).
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

func CowCreateService(r *cow.CreateRequest) error {
	item := model.Cow{
		Name:       r.Name,
		Tag:        r.Tag,
		Age:        r.Age,
		CanMilking: r.CanMilking,
		Status:     r.Status,
		Condition:  r.Condition,
	}

	return pg.DB.Create(&item).Error
}

func CowUpdateService(r *cow.UpdateRequest) error {
	var item model.Cow
	err := pg.DB.Where("id = ?", r.ID).First(&item).Error
	if err != nil {
		return err
	}

	item.Name = r.Name
	item.Tag = r.Tag
	item.Age = r.Age
	item.CanMilking = r.CanMilking
	item.Status = r.Status
	item.Condition = r.Condition

	return pg.DB.Save(&item).Error
}
