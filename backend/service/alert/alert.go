package alert

import (
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	alertdto "github.com/zcl0621/compx576-smart-dairy-system/dto/alert"
	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"gorm.io/gorm"
)

func ListService(r *alertdto.ListQuery) (*alertdto.ListResponse, error) {
	db := alertBaseQuery()
	db = applyAlertFilters(db, r.CowID, r.Severity)
	db = db.Order("alerts.created_at desc")

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	var list []alertdto.ListItem
	if err := db.Offset(r.GetOffset()).Limit(r.GetLimit()).Find(&list).Error; err != nil {
		return nil, err
	}

	response := alertdto.ListResponse{
		List: list,
		PageResponse: common.PageResponse{
			Page:       r.GetPage(),
			Total:      total,
			TotalPages: r.GetTotalPages(total),
		},
	}

	return &response, nil
}

func SummaryService() (*alertdto.SummaryResponse, error) {
	response := alertdto.SummaryResponse{}

	active, err := countActiveAlerts()
	if err != nil {
		return nil, err
	}
	response.Active = active

	warning, err := countActiveAlerts(model.AlertSeverityWarning)
	if err != nil {
		return nil, err
	}
	response.Warning = warning

	critical, err := countActiveAlerts(model.AlertSeverityCritical)
	if err != nil {
		return nil, err
	}
	response.Critical = critical

	offline, err := countActiveAlerts(model.AlertSeverityOffline)
	if err != nil {
		return nil, err
	}
	response.Offline = offline

	return &response, nil
}

func alertBaseQuery() *gorm.DB {
	return pg.DB.Model(&model.Alert{}).
		Select(
			"alerts.id",
			"alerts.created_at",
			"alerts.updated_at",
			"alerts.cow_id",
			"(SELECT name FROM cows WHERE cows.id = alerts.cow_id) AS cow_name",
			"alerts.report_id",
			"alerts.metric_key",
			"alerts.title",
			"alerts.message",
			"alerts.severity",
			"alerts.status",
			"alerts.resolved_at",
		)
}

func applyAlertFilters(db *gorm.DB, cowID string, severity string) *gorm.DB {
	db = db.Where("status = ?", model.AlertStatusActive)
	if cowID != "" {
		db = db.Where("cow_id = ?", cowID)
	}
	if severity != "" {
		db = db.Where("severity = ?", severity)
	}
	return db
}

func activeAlertSummaryQuery() *gorm.DB {
	return pg.DB.Model(&model.Alert{}).Where("status = ?", model.AlertStatusActive)
}

func countActiveAlerts(severities ...model.AlertSeverity) (int64, error) {
	db := activeAlertSummaryQuery()
	if len(severities) > 0 {
		db = db.Where("severity = ?", severities[0])
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
