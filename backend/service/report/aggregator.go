package report

import (
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"gorm.io/gorm"
)

// expectedHealthReadings is the expected reading count for one health metric
// over 7 days at 30s intervals.
const expectedHealthReadings = 7 * 24 * 60 * 2 // 20160

// AggregateCowMetrics fetches stats for the current window and the previous 7-day baseline.
// Both windows are needed so score can compare weight/milk against the cow's own history.
func AggregateCowMetrics(
	db *gorm.DB,
	cowID string,
	currentStart, currentEnd, baselineStart, baselineEnd time.Time,
) (*CowStats, error) {
	current, err := windowStats(db, cowID, currentStart, currentEnd)
	if err != nil {
		return nil, err
	}
	baseline, err := windowStats(db, cowID, baselineStart, baselineEnd)
	if err != nil {
		return nil, err
	}
	return &CowStats{CowID: cowID, Current: current, Baseline: baseline}, nil
}

// windowStats fetches all five metrics in one GROUP BY query.
// Abnormal thresholds are embedded as CASE expressions — one round trip per window.
func windowStats(db *gorm.DB, cowID string, start, end time.Time) (WindowStats, error) {
	var rows []struct {
		MetricType string
		Count      int64
		Avg        *float64
		Min        *float64
		Max        *float64
		Abnormal   int64
	}
	err := db.Raw(`
        SELECT metric_type,
               COUNT(*) AS count,
               AVG(metric_value) AS avg,
               MIN(metric_value) AS min,
               MAX(metric_value) AS max,
               SUM(CASE
                   WHEN metric_type = 'temperature'   AND (metric_value < 38.0 OR metric_value > 39.0) THEN 1
                   WHEN metric_type = 'heart_rate'    AND (metric_value < 48   OR metric_value > 84)   THEN 1
                   WHEN metric_type = 'blood_oxygen'  AND  metric_value < 90                           THEN 1
                   ELSE 0
               END) AS abnormal
        FROM metrics
        WHERE cow_id = ? AND metric_type IN (?,?,?,?,?)
          AND created_at >= ? AND created_at < ?
          AND deleted_at IS NULL
        GROUP BY metric_type`,
		cowID,
		model.MetricTypeTemperature, model.MetricTypeHeartRate, model.MetricTypeBloodOxygen,
		model.MetricTypeWeight, model.MetricTypeMilkAmount,
		start, end,
	).Scan(&rows).Error
	if err != nil {
		return WindowStats{}, err
	}

	var w WindowStats
	for _, r := range rows {
		s := MetricStats{
			Count:         int(r.Count),
			Avg:           *r.Avg,
			Min:           *r.Min,
			Max:           *r.Max,
			AbnormalCount: int(r.Abnormal),
		}
		switch model.MetricType(r.MetricType) {
		case model.MetricTypeTemperature:
			w.Temperature = s
		case model.MetricTypeHeartRate:
			w.HeartRate = s
		case model.MetricTypeBloodOxygen:
			w.BloodOxygen = s
		case model.MetricTypeWeight:
			w.Weight = s
		case model.MetricTypeMilkAmount:
			w.MilkAmount = s
		}
	}

	// Coverage averages the three health-metric coverage ratios (each capped at 100).
	cov := func(count int) float64 {
		ratio := float64(count) / float64(expectedHealthReadings) * 100.0
		if ratio > 100 {
			return 100
		}
		return ratio
	}
	w.CoveragePct = (cov(w.Temperature.Count) + cov(w.HeartRate.Count) + cov(w.BloodOxygen.Count)) / 3.0
	return w, nil
}
