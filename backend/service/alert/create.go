package alert

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/clause"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/mail"
	userservice "github.com/zcl0621/compx576-smart-dairy-system/service/user"
)

// CreateIfNotExists skips silently on duplicate — deduplication is enforced by the
// partial unique index on (cow_id, metric_key) WHERE status='active'.
func CreateIfNotExists(cowID string, metricKey model.MetricType, severity model.AlertSeverity, title, message string) error {
	a := &model.Alert{
		CowID:     cowID,
		MetricKey: metricKey,
		Severity:  severity,
		Status:    model.AlertStatusActive,
		Title:     title,
		Message:   message,
	}
	result := pg.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(a)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}

	var cow model.Cow
	if err := pg.DB.Select("name").First(&cow, "id = ?", cowID).Error; err != nil {
		projectlog.L().Warn("alert email: cow lookup failed", zap.Error(err))
		return nil
	}
	emails, err := userservice.GetAllEmails()
	if err != nil || len(emails) == 0 {
		projectlog.L().Warn("alert email: no recipients", zap.Error(err))
		return nil
	}
	if err := mail.SendAlertEmail(emails, cow.Name, title, string(severity), message); err != nil {
		projectlog.L().Warn("alert email: send failed", zap.Error(err))
	}
	return nil
}

func ResolveIfExists(cowID string, metricKey model.MetricType) error {
	return pg.DB.Model(&model.Alert{}).
		Where("cow_id = ? AND metric_key = ? AND status = ?", cowID, metricKey, model.AlertStatusActive).
		Updates(map[string]any{
			"status":      model.AlertStatusResolved,
			"resolved_at": time.Now(),
		}).Error
}
