package metric

import (
	"errors"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	cowdto "github.com/zcl0621/compx576-smart-dairy-system/dto/cow"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrBadMetricRange = errors.New("metric range is wrong")

type metricRow struct {
	CreatedAt   time.Time `gorm:"column:created_at"`
	MetricValue float64   `gorm:"column:metric_value"`
}

type coordinateMetricRow struct {
	CreatedAt   time.Time        `gorm:"column:created_at"`
	MetricType  model.MetricType `gorm:"column:metric_type"`
	MetricValue float64          `gorm:"column:metric_value"`
}

type movementPointRow struct {
	CreatedAt time.Time
	Latitude  float64
	Longitude float64
}

const movementPairTolerance = 5 * time.Minute

func TemperatureService(r *cowdto.MetricQuery) (*cowdto.TemperatureMetricResponse, error) {
	metricRange, err := parseMetricRange(r.CowID, r.Range)
	if err != nil {
		return nil, err
	}

	rows, err := loadMetricRows(r.CowID, model.MetricTypeTemperature, metricRange)
	if err != nil {
		return nil, err
	}

	response := &cowdto.TemperatureMetricResponse{
		CowID: r.CowID,
		Range: metricRange,
		Summary: cowdto.TemperatureMetricSummary{
			Status: statusForTemperature(rows),
		},
		Series: toMetricPoints(rows),
	}
	fillMetricStats(rows, &response.UpdatedAt, &response.Summary.Current, &response.Summary.Avg, &response.Summary.Min, &response.Summary.Max)
	logMetricSummary("temperature", r.CowID, metricRange, len(rows), string(response.Summary.Status), len(response.Series))
	return response, nil
}

func HeartRateService(r *cowdto.MetricQuery) (*cowdto.HeartRateMetricResponse, error) {
	metricRange, err := parseMetricRange(r.CowID, r.Range)
	if err != nil {
		return nil, err
	}

	rows, err := loadMetricRows(r.CowID, model.MetricTypeHeartRate, metricRange)
	if err != nil {
		return nil, err
	}

	response := &cowdto.HeartRateMetricResponse{
		CowID: r.CowID,
		Range: metricRange,
		Summary: cowdto.HeartRateMetricSummary{
			Status: statusForHeartRate(rows),
		},
		Series: toMetricPoints(rows),
	}
	fillMetricStats(rows, &response.UpdatedAt, &response.Summary.Current, &response.Summary.Avg, &response.Summary.Min, &response.Summary.Max)
	logMetricSummary("heart rate", r.CowID, metricRange, len(rows), string(response.Summary.Status), len(response.Series))
	return response, nil
}

func BloodOxygenService(r *cowdto.MetricQuery) (*cowdto.BloodOxygenMetricResponse, error) {
	metricRange, err := parseMetricRange(r.CowID, r.Range)
	if err != nil {
		return nil, err
	}

	rows, err := loadMetricRows(r.CowID, model.MetricTypeBloodOxygen, metricRange)
	if err != nil {
		return nil, err
	}

	response := &cowdto.BloodOxygenMetricResponse{
		CowID: r.CowID,
		Range: metricRange,
		Summary: cowdto.BloodOxygenMetricSummary{
			Status: statusForBloodOxygen(rows),
		},
		Series: toMetricPoints(rows),
	}
	fillMetricStats(rows, &response.UpdatedAt, &response.Summary.Current, &response.Summary.Avg, &response.Summary.Min, &response.Summary.Max)
	logMetricSummary("blood oxygen", r.CowID, metricRange, len(rows), string(response.Summary.Status), len(response.Series))
	return response, nil
}

func WeightService(r *cowdto.MetricQuery) (*cowdto.WeightMetricResponse, error) {
	metricRange, err := parseMetricRange(r.CowID, r.Range)
	if err != nil {
		return nil, err
	}

	rows, err := loadMetricRows(r.CowID, model.MetricTypeWeight, metricRange)
	if err != nil {
		return nil, err
	}

	response := &cowdto.WeightMetricResponse{
		CowID:  r.CowID,
		Range:  metricRange,
		Series: toMetricPoints(rows),
	}
	fillMetricStats(rows, &response.UpdatedAt, &response.Summary.Current, &response.Summary.Avg, &response.Summary.Min, &response.Summary.Max)
	logMetricSummary("weight", r.CowID, metricRange, len(rows), "n/a", len(response.Series))
	return response, nil
}

func MilkAmountService(r *cowdto.MetricQuery) (*cowdto.MilkAmountMetricResponse, error) {
	metricRange, err := parseMetricRange(r.CowID, r.Range)
	if err != nil {
		return nil, err
	}

	rows, err := loadMetricRows(r.CowID, model.MetricTypeMilkAmount, metricRange)
	if err != nil {
		return nil, err
	}

	response := &cowdto.MilkAmountMetricResponse{
		CowID:  r.CowID,
		Range:  metricRange,
		Series: toMetricPoints(rows),
	}
	if len(rows) > 0 {
		updatedAt := rows[len(rows)-1].CreatedAt
		response.UpdatedAt = &updatedAt
	}
	response.Summary.SessionCount = int64(len(rows))
	for _, row := range rows {
		response.Summary.Total += row.MetricValue
	}
	if response.Summary.SessionCount > 0 {
		response.Summary.AvgPerSession = response.Summary.Total / float64(response.Summary.SessionCount)
	}
	projectlog.L().Debug("build milk summary",
		zap.String("cow_id", r.CowID),
		zap.String("range", string(metricRange)),
		zap.Int("rows", len(rows)),
		zap.Int("series", len(response.Series)),
		zap.Int64("session_count", response.Summary.SessionCount),
	)
	return response, nil
}

func MovementService(r *cowdto.MetricQuery) (*cowdto.MovementMetricResponse, error) {
	metricRange, err := parseMetricRange(r.CowID, r.Range)
	if err != nil {
		return nil, err
	}

	rows, err := loadMovementRows(r.CowID, metricRange)
	if err != nil {
		return nil, err
	}

	points := buildMovementPoints(r.CowID, rows)

	response := &cowdto.MovementMetricResponse{
		CowID: r.CowID,
		Range: metricRange,
		Summary: cowdto.MovementMetricSummary{
			PointCount: int64(len(points)),
			Status:     statusForMovement(points),
		},
		Series: make([]cowdto.MovementPoint, 0),
	}
	if len(points) < 2 {
		projectlog.L().Debug("movement rows too short",
			zap.String("cow_id", r.CowID),
			zap.String("range", string(metricRange)),
			zap.Int("points", len(points)),
		)
	}
	if len(points) > 0 {
		updatedAt := points[len(points)-1].CreatedAt
		response.UpdatedAt = &updatedAt
	}

	for i := 1; i < len(points); i++ {
		prev := points[i-1]
		curr := points[i]
		distance := util.HaversineMeters(prev.Latitude, prev.Longitude, curr.Latitude, curr.Longitude)
		response.Summary.DistanceM += distance
		response.Series = append(response.Series, cowdto.MovementPoint{Time: curr.CreatedAt, DistanceM: distance})
	}

	projectlog.L().Debug("build movement summary",
		zap.String("cow_id", r.CowID),
		zap.String("range", string(metricRange)),
		zap.Int("rows", len(rows)),
		zap.Int("points", len(points)),
		zap.Int("series", len(response.Series)),
		zap.String("status", string(response.Summary.Status)),
	)

	return response, nil
}

func parseMetricRange(cowID string, raw string) (model.MetricRange, error) {
	metricRange := model.MetricRange(raw)
	switch metricRange {
	case "", model.MetricRange24H:
		return model.MetricRange24H, nil
	case model.MetricRange7D, model.MetricRange30D, model.MetricRangeAll:
		return metricRange, nil
	default:
		projectlog.L().Warn("metric range is wrong",
			zap.String("cow_id", cowID),
			zap.String("range", raw),
		)
		return "", ErrBadMetricRange
	}
}

func loadMetricRows(cowID string, metricType model.MetricType, metricRange model.MetricRange) ([]metricRow, error) {
	db := applyMetricRange(pg.DB.Model(&model.Metric{}), metricRange).
		Select("created_at, metric_value").
		Where("cow_id = ? AND metric_type = ?", cowID, metricType).
		Order("created_at asc")

	var rows []metricRow
	if err := db.Find(&rows).Error; err != nil {
		projectlog.L().Error("load metric rows failed",
			zap.String("cow_id", cowID),
			zap.String("metric_type", string(metricType)),
			zap.String("range", string(metricRange)),
			zap.Error(err),
		)
		return nil, err
	}
	projectlog.L().Debug("loaded metric rows",
		zap.String("cow_id", cowID),
		zap.String("metric_type", string(metricType)),
		zap.String("range", string(metricRange)),
		zap.Int("rows", len(rows)),
	)
	return rows, nil
}

func loadMovementRows(cowID string, metricRange model.MetricRange) ([]coordinateMetricRow, error) {
	db := applyMetricRange(pg.DB.Model(&model.Metric{}), metricRange).
		Select("created_at, metric_type, metric_value").
		Where("cow_id = ? AND metric_type IN (?, ?)", cowID, model.MetricTypeLatitude, model.MetricTypeLongitude).
		Order("created_at asc, metric_type asc")

	var rows []coordinateMetricRow
	if err := db.Find(&rows).Error; err != nil {
		projectlog.L().Error("load movement rows failed",
			zap.String("cow_id", cowID),
			zap.String("range", string(metricRange)),
			zap.Error(err),
		)
		return nil, err
	}
	projectlog.L().Debug("loaded movement rows",
		zap.String("cow_id", cowID),
		zap.String("range", string(metricRange)),
		zap.Int("rows", len(rows)),
	)
	return rows, nil
}

func buildMovementPoints(cowID string, rows []coordinateMetricRow) []movementPointRow {
	points := make([]movementPointRow, 0)
	var latitudeRow *coordinateMetricRow
	var longitudeRow *coordinateMetricRow

	flushPair := func() {
		if latitudeRow == nil || longitudeRow == nil {
			return
		}

		pairTime := latitudeRow.CreatedAt
		if longitudeRow.CreatedAt.After(pairTime) {
			pairTime = longitudeRow.CreatedAt
		}

		points = append(points, movementPointRow{
			CreatedAt: pairTime,
			Latitude:  latitudeRow.MetricValue,
			Longitude: longitudeRow.MetricValue,
		})
		latitudeRow = nil
		longitudeRow = nil
	}

	for i := range rows {
		row := &rows[i]
		switch row.MetricType {
		case model.MetricTypeLatitude:
			if latitudeRow != nil {
				projectlog.L().Debug("skip movement point missing lat or lng",
					zap.String("cow_id", cowID),
					zap.Time("time", latitudeRow.CreatedAt),
				)
			}
			if longitudeRow != nil && row.CreatedAt.Sub(longitudeRow.CreatedAt) > movementPairTolerance {
				projectlog.L().Debug("skip movement point missing lat or lng",
					zap.String("cow_id", cowID),
					zap.Time("time", longitudeRow.CreatedAt),
				)
				longitudeRow = nil
			}
			latitudeRow = row
		case model.MetricTypeLongitude:
			if longitudeRow != nil {
				projectlog.L().Debug("skip movement point missing lat or lng",
					zap.String("cow_id", cowID),
					zap.Time("time", longitudeRow.CreatedAt),
				)
			}
			if latitudeRow != nil && row.CreatedAt.Sub(latitudeRow.CreatedAt) > movementPairTolerance {
				projectlog.L().Debug("skip movement point missing lat or lng",
					zap.String("cow_id", cowID),
					zap.Time("time", latitudeRow.CreatedAt),
				)
				latitudeRow = nil
			}
			longitudeRow = row
		}

		if latitudeRow != nil && longitudeRow != nil {
			diff := latitudeRow.CreatedAt.Sub(longitudeRow.CreatedAt)
			if diff < 0 {
				diff = -diff
			}
			if diff <= movementPairTolerance {
				flushPair()
			}
		}
	}

	if latitudeRow != nil {
		projectlog.L().Debug("skip movement point missing lat or lng",
			zap.String("cow_id", cowID),
			zap.Time("time", latitudeRow.CreatedAt),
		)
	}
	if longitudeRow != nil {
		projectlog.L().Debug("skip movement point missing lat or lng",
			zap.String("cow_id", cowID),
			zap.Time("time", longitudeRow.CreatedAt),
		)
	}

	return points
}

func applyMetricRange(db *gorm.DB, metricRange model.MetricRange) *gorm.DB {
	if metricRange == model.MetricRangeAll {
		return db
	}

	now := time.Now()
	switch metricRange {
	case model.MetricRange24H:
		return db.Where("created_at >= ?", now.Add(-24*time.Hour))
	case model.MetricRange7D:
		return db.Where("created_at >= ?", now.AddDate(0, 0, -7))
	case model.MetricRange30D:
		return db.Where("created_at >= ?", now.AddDate(0, 0, -30))
	default:
		return db
	}
}

func fillMetricStats(rows []metricRow, updatedAt **time.Time, current **float64, avg **float64, min **float64, max **float64) {
	if len(rows) == 0 {
		return
	}

	latest := rows[len(rows)-1]
	*updatedAt = &latest.CreatedAt
	*current = util.Float64Ptr(latest.MetricValue)

	total := 0.0
	minValue := rows[0].MetricValue
	maxValue := rows[0].MetricValue
	for _, row := range rows {
		total += row.MetricValue
		if row.MetricValue < minValue {
			minValue = row.MetricValue
		}
		if row.MetricValue > maxValue {
			maxValue = row.MetricValue
		}
	}
	avgValue := total / float64(len(rows))
	*avg = util.Float64Ptr(avgValue)
	*min = util.Float64Ptr(minValue)
	*max = util.Float64Ptr(maxValue)
}

func toMetricPoints(rows []metricRow) []cowdto.MetricPoint {
	points := make([]cowdto.MetricPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, cowdto.MetricPoint{Time: row.CreatedAt, Value: row.MetricValue})
	}
	return points
}

func statusForTemperature(rows []metricRow) model.ReportMetricStatus {
	if len(rows) == 0 {
		return model.ReportMetricStatusOffline
	}
	current := rows[len(rows)-1].MetricValue
	switch {
	case current >= 40:
		return model.ReportMetricStatusCritical
	case current >= 39:
		return model.ReportMetricStatusWarning
	default:
		return model.ReportMetricStatusNormal
	}
}

func statusForHeartRate(rows []metricRow) model.ReportMetricStatus {
	if len(rows) == 0 {
		return model.ReportMetricStatusOffline
	}
	current := rows[len(rows)-1].MetricValue
	switch {
	case current >= 100 || current <= 45:
		return model.ReportMetricStatusCritical
	case current >= 90 || current <= 50:
		return model.ReportMetricStatusWarning
	default:
		return model.ReportMetricStatusNormal
	}
}

func statusForBloodOxygen(rows []metricRow) model.ReportMetricStatus {
	if len(rows) == 0 {
		return model.ReportMetricStatusOffline
	}
	current := rows[len(rows)-1].MetricValue
	switch {
	case current < 90:
		return model.ReportMetricStatusCritical
	case current < 95:
		return model.ReportMetricStatusWarning
	default:
		return model.ReportMetricStatusNormal
	}
}

func statusForMovement(points []movementPointRow) model.ReportMetricStatus {
	if len(points) == 0 {
		return model.ReportMetricStatusOffline
	}
	if len(points) < 2 {
		return model.ReportMetricStatusWarning
	}
	return model.ReportMetricStatusNormal
}

func logMetricSummary(metricName string, cowID string, metricRange model.MetricRange, rowCount int, status string, seriesCount int) {
	projectlog.L().Debug("build metric summary",
		zap.String("metric", metricName),
		zap.String("cow_id", cowID),
		zap.String("range", string(metricRange)),
		zap.Int("rows", rowCount),
		zap.Int("series", seriesCount),
		zap.String("status", status),
	)
}
