package alert

import (
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

// sevenDayAverage returns the average and count of metric_value over the 7 days before `before`.
func sevenDayAverage(cowID string, metricType model.MetricType, before time.Time) (avg float64, count int64, err error) {
	sevenDaysAgo := before.AddDate(0, 0, -7)
	var result struct {
		Avg   *float64
		Count int64
	}
	err = pg.DB.Raw(
		`SELECT AVG(metric_value) AS avg, COUNT(*) AS count FROM metrics
		 WHERE cow_id = ? AND metric_type = ? AND created_at < ? AND created_at >= ? AND deleted_at IS NULL`,
		cowID, metricType, before, sevenDaysAgo,
	).Scan(&result).Error
	if err != nil || result.Count == 0 {
		return 0, result.Count, err
	}
	return *result.Avg, result.Count, nil
}
