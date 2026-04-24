package alert

import (
	"fmt"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

type healthReading struct {
	MetricValue float64
}

// EvaluateHealth evaluates the 3 most recent readings for cowID + metricType.
// Resolve runs independently of the 3-reading threshold so a recovering cow
// isn't stuck waiting for 3 new normal readings before the alert clears.
func EvaluateHealth(cowID string, metricType model.MetricType) error {
	var rows []healthReading
	if err := pg.DB.Model(&model.Metric{}).
		Select("metric_value").
		Where("cow_id = ? AND metric_type = ?", cowID, metricType).
		Order("created_at DESC").
		Limit(3).
		Find(&rows).Error; err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}

	if isHealthNormal(metricType, rows[0].MetricValue) {
		return ResolveIfExists(cowID, metricType)
	}

	if len(rows) < 3 {
		return nil
	}
	for _, r := range rows {
		if isHealthNormal(metricType, r.MetricValue) {
			return nil
		}
	}

	severity := healthSeverity(metricType, rows)
	title, message := healthAlertText(metricType, rows[0].MetricValue)
	return CreateIfNotExists(cowID, metricType, severity, title, message)
}

func isHealthNormal(metricType model.MetricType, v float64) bool {
	return healthSeverityForValue(metricType, v) == ""
}

func healthSeverityForValue(metricType model.MetricType, v float64) model.AlertSeverity {
	switch metricType {
	case model.MetricTypeTemperature:
		switch {
		case v > TempCritHigh || v < TempCritLow:
			return model.AlertSeverityCritical
		case v > TempWarnHigh || v < TempWarnLow:
			return model.AlertSeverityWarning
		}
	case model.MetricTypeHeartRate:
		switch {
		case v > HRCritHigh || v < HRCritLow:
			return model.AlertSeverityCritical
		case v > HRWarnHigh || v < HRWarnLow:
			return model.AlertSeverityWarning
		}
	case model.MetricTypeBloodOxygen:
		switch {
		case v < BOCrit:
			return model.AlertSeverityCritical
		case v < BOWarn:
			return model.AlertSeverityWarning
		}
	}
	return ""
}

func healthSeverity(metricType model.MetricType, rows []healthReading) model.AlertSeverity {
	for _, r := range rows {
		if healthSeverityForValue(metricType, r.MetricValue) == model.AlertSeverityCritical {
			return model.AlertSeverityCritical
		}
	}
	return model.AlertSeverityWarning
}

func healthAlertText(metricType model.MetricType, current float64) (title, message string) {
	switch metricType {
	case model.MetricTypeTemperature:
		title = "Abnormal temperature"
		message = fmt.Sprintf("Temperature %.1f°C is outside the normal range (38.0–39.5°C)", current)
	case model.MetricTypeHeartRate:
		title = "Abnormal heart rate"
		message = fmt.Sprintf("Heart rate %.0f bpm is outside the normal range (48–84 bpm)", current)
	case model.MetricTypeBloodOxygen:
		title = "Low blood oxygen"
		message = fmt.Sprintf("Blood oxygen %.1f%% is below the normal threshold (≥95%%)", current)
	default:
		title = "Abnormal reading"
		message = fmt.Sprintf("Value %.2f is outside the normal range", current)
	}
	return
}
