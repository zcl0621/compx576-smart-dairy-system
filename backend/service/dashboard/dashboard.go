package dashboard

import (
	"fmt"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	dashboarddto "github.com/zcl0621/compx576-smart-dairy-system/dto/dashboard"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"gorm.io/gorm"
)

func SummaryService() (*dashboarddto.SummaryResponse, error) {
	response := dashboarddto.SummaryResponse{}

	total, err := countDashboardCows()
	if err != nil {
		return nil, err
	}
	response.TotalCows = total

	normal, err := countDashboardCows(model.CowConditionNormal)
	if err != nil {
		return nil, err
	}
	response.Normal = normal

	warning, err := countDashboardCows(model.CowConditionWarning)
	if err != nil {
		return nil, err
	}
	response.Warning = warning

	critical, err := countDashboardCows(model.CowConditionCritical)
	if err != nil {
		return nil, err
	}
	response.Critical = critical

	offline, err := countDashboardCows(model.CowConditionOffline)
	if err != nil {
		return nil, err
	}
	response.Offline = offline

	return &response, nil
}

func ListService(r *dashboarddto.ListQuery) (*dashboarddto.ListResponse, error) {
	var total int64
	countDB := dashboardCowBaseQuery()
	if err := countDB.Count(&total).Error; err != nil {
		return nil, err
	}

	db := dashboardListBaseQuery()
	db = model.OrderByCowCondition(db)
	db = db.Order("updated_at desc")

	var list []dashboarddto.ListItem
	if err := db.Offset(r.GetOffset()).Limit(r.GetLimit()).Find(&list).Error; err != nil {
		return nil, err
	}

	response := dashboarddto.ListResponse{
		List: list,
		PageResponse: common.PageResponse{
			Page:       r.GetPage(),
			Total:      total,
			TotalPages: r.GetTotalPages(total),
		},
	}

	return &response, nil
}

func dashboardCowBaseQuery() *gorm.DB {
	return pg.DB.Model(&model.Cow{}).
		Where("status != ?", model.CowStatusSold)
}

func dashboardListBaseQuery() *gorm.DB {
	return dashboardCowBaseQuery().
		Select(
			"cows.id",
			"cows.name",
			"cows.tag",
			model.ConditionSubQuery()+" AS condition",
			"cows.updated_at",
			latestMetricValueSubQuery(model.MetricTypeTemperature)+" AS temperature",
			latestMetricValueSubQuery(model.MetricTypeHeartRate)+" AS heart_rate",
			latestMetricValueSubQuery(model.MetricTypeBloodOxygen)+" AS blood_oxygen",
			latestAlertMessageSubQuery()+" AS alert_message",
		)
}

func latestMetricValueSubQuery(metricType model.MetricType) string {
	return fmt.Sprintf(`(
		SELECT metric_value
		FROM metrics
		WHERE metrics.cow_id = cows.id
		  AND metrics.metric_type = '%s'
		ORDER BY metrics.created_at DESC
		LIMIT 1
	)`, metricType)
}

func latestAlertMessageSubQuery() string {
	return fmt.Sprintf(`(
		SELECT message
		FROM alerts
		WHERE alerts.cow_id = cows.id
		  AND alerts.status = '%s'
		ORDER BY CASE alerts.severity
			WHEN '%s' THEN 1
			WHEN '%s' THEN 2
			WHEN '%s' THEN 3
			ELSE 999
		END,
		alerts.created_at DESC
		LIMIT 1
	)`, model.AlertStatusActive, model.AlertSeverityCritical, model.AlertSeverityWarning, model.AlertSeverityOffline)
}

func countDashboardCows(conditions ...model.CowCondition) (int64, error) {
	db := dashboardCowBaseQuery()
	if len(conditions) > 0 {
		switch conditions[0] {
		case model.CowConditionNormal:
			db = db.Where("NOT EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = 'active' AND alerts.deleted_at IS NULL)")
		default:
			db = db.Where("EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = 'active' AND alerts.severity = ? AND alerts.deleted_at IS NULL)", conditions[0])
		}
	}
	var total int64
	return total, db.Count(&total).Error
}
