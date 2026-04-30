package report

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	alertservice "github.com/zcl0621/compx576-smart-dairy-system/service/alert"
	"github.com/zcl0621/compx576-smart-dairy-system/service/llm"
)

// advisoryLockID is the pg_try_advisory_lock key for the daily report job.
// Arbitrary number — just needs to be unique within this codebase.
const advisoryLockID int64 = 9176

type LLM interface {
	Generate(ctx context.Context, system, user string) (llm.LLMOutput, error)
}

type Generator struct {
	DB  *gorm.DB
	LLM LLM
	Now func() time.Time
	// SleepFn lets tests skip the retry backoff. Defaults to time.After.
	SleepFn func(ctx context.Context, d time.Duration)
}

func AcquireLock(ctx context.Context, db *gorm.DB) (*sql.Conn, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("acquire lock: get sql.DB: %w", err)
	}
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquire lock: open conn: %w", err)
	}

	var got bool
	row := conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", advisoryLockID)
	if err := row.Scan(&got); err != nil {
		conn.Close()
		return nil, fmt.Errorf("acquire lock: scan: %w", err)
	}
	if !got {
		conn.Close()
		return nil, nil
	}
	return conn, nil
}

func DeleteTodaysReports(db *gorm.DB) error {
	return db.Exec(`DELETE FROM reports
        WHERE (created_at AT TIME ZONE 'Pacific/Auckland')::date
            = (now() AT TIME ZONE 'Pacific/Auckland')::date`).Error
}

func ListActiveCows(db *gorm.DB) ([]string, error) {
	var ids []string
	err := db.Model(&model.Cow{}).
		Where("status = ? AND deleted_at IS NULL", model.CowStatusInFarm).
		Order("id").
		Pluck("id", &ids).Error
	return ids, err
}

// nzLoc loads once at init. Missing tzdata breaks the job anyway,
// so panicking here surfaces the problem early.
var nzLoc = func() *time.Location {
	loc, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		panic(fmt.Sprintf("report: load Pacific/Auckland: %v", err))
	}
	return loc
}()

// PeriodWindow returns the 7-day window ending at midnight NZ time today.
func PeriodWindow(now time.Time) (start, end time.Time) {
	end = time.Date(now.In(nzLoc).Year(), now.In(nzLoc).Month(), now.In(nzLoc).Day(),
		0, 0, 0, 0, nzLoc).UTC()
	start = end.AddDate(0, 0, -7)
	return start, end
}

// RunOnce is the entry point for cmd/report. Returns nil on success and when
// another instance holds the lock (logs and exits 0 in that case).
func (g *Generator) RunOnce(ctx context.Context) error {
	conn, err := AcquireLock(ctx, g.DB)
	if err != nil {
		return err
	}
	if conn == nil {
		projectlog.L().Info("report generator: another instance is running, exiting")
		return nil
	}
	defer conn.Close()

	if err := DeleteTodaysReports(g.DB); err != nil {
		return fmt.Errorf("delete today: %w", err)
	}

	cowIDs, err := ListActiveCows(g.DB)
	if err != nil {
		return fmt.Errorf("list cows: %w", err)
	}
	projectlog.L().Info("report generator started", zap.Int("cows", len(cowIDs)))

	sched := NewRetryScheduler(cowIDs)
	if g.SleepFn != nil {
		sched.sleep = g.SleepFn
	}
	sched.SetFinalFailureHook(func(ctx context.Context, cowID string) {
		if err := g.recordFinalFailure(cowID); err != nil {
			projectlog.L().Warn("record final failure", zap.String("cow_id", cowID), zap.Error(err))
		}
	})

	return sched.Run(ctx, func(ctx context.Context, cowID string) error {
		if err := g.GenerateOne(ctx, cowID); err != nil {
			projectlog.L().Warn("generate one failed",
				zap.String("cow_id", cowID), zap.Error(err))
			return err
		}
		// resolve any prior failure alert on success
		if rerr := alertservice.ResolveIfExists(cowID, model.MetricTypeReportFailure); rerr != nil {
			projectlog.L().Warn("resolve prior failure alert",
				zap.String("cow_id", cowID), zap.Error(rerr))
		}
		return nil
	})
}

// GenerateOne runs the full per-cow report flow.
func (g *Generator) GenerateOne(ctx context.Context, cowID string) error {
	var cow model.Cow
	if err := g.DB.Select("id", "name").First(&cow, "id = ?", cowID).Error; err != nil {
		return fmt.Errorf("load cow: %w", err)
	}

	now := g.Now()
	start, end := PeriodWindow(now)
	baselineEnd := start
	baselineStart := baselineEnd.AddDate(0, 0, -7)

	stats, err := AggregateCowMetrics(g.DB, cowID, start, end, baselineStart, baselineEnd)
	if err != nil {
		return fmt.Errorf("aggregate: %w", err)
	}

	// cap at 20: alert deduction saturates around 10 and the details column stays small
	var alerts []model.Alert
	if err := g.DB.Where("cow_id = ? AND created_at >= ? AND created_at < ?",
		cowID, start, end).
		Order("created_at asc").
		Limit(20).
		Find(&alerts).Error; err != nil {
		return fmt.Errorf("load alerts: %w", err)
	}

	score := CalculateScore(stats, alerts)
	out, err := g.LLM.Generate(ctx, BuildSystemPrompt(), BuildUserPrompt(&cow, stats, alerts, score, start, end))
	if err != nil {
		return fmt.Errorf("llm: %w", err)
	}

	details := buildDetails(stats, alerts, out.Note)
	rep := &model.Report{
		CowID:       cowID,
		PeriodStart: start,
		PeriodEnd:   end,
		Summary:     out.Summary,
		Score:       score,
		Details:     details,
	}
	if err := g.DB.Create(rep).Error; err != nil {
		return fmt.Errorf("insert report: %w", err)
	}
	return nil
}

func (g *Generator) recordFinalFailure(cowID string) error {
	var cow model.Cow
	if err := g.DB.Select("name").First(&cow, "id = ?", cowID).Error; err != nil {
		return fmt.Errorf("lookup cow: %w", err)
	}
	return alertservice.CreateIfNotExists(
		cowID,
		model.MetricTypeReportFailure,
		model.AlertSeverityCritical,
		"Report generation failed",
		fmt.Sprintf("Report generation failed for %s. Manual check required.", cow.Name),
	)
}

func buildDetails(stats *CowStats, alerts []model.Alert, note string) model.ReportDetails {
	var metrics []model.ReportMetric
	types := []struct {
		key   model.MetricType
		label string
		unit  model.MetricUnit
		stats MetricStats
		base  float64
	}{
		{model.MetricTypeTemperature, "Temperature", model.MetricUnitCelsius, stats.Current.Temperature, 0},
		{model.MetricTypeHeartRate, "Heart Rate", model.MetricUnitBPM, stats.Current.HeartRate, 0},
		{model.MetricTypeBloodOxygen, "Blood Oxygen", model.MetricUnitPercent, stats.Current.BloodOxygen, 0},
		{model.MetricTypeWeight, "Weight", model.MetricUnitKG, stats.Current.Weight, stats.Baseline.Weight.Avg},
		{model.MetricTypeMilkAmount, "Milk Yield", model.MetricUnitLiters, stats.Current.MilkAmount, stats.Baseline.MilkAmount.Avg},
	}
	for _, m := range types {
		var status model.ReportMetricStatus
		var high bool
		if m.stats.Count == 0 {
			status = model.ReportMetricStatusOffline
		} else {
			lv := ClassifyDeviation(m.key, m.stats.Avg, m.base)
			status = DeviationToStatus(lv)
			switch m.key {
			case model.MetricTypeTemperature:
				high = m.stats.Avg > 39.0
			case model.MetricTypeHeartRate:
				high = m.stats.Avg > 84
			case model.MetricTypeBloodOxygen:
				high = false
			default:
				high = m.base != 0 && m.stats.Avg > m.base
			}
		}
		metrics = append(metrics, model.ReportMetric{
			Key:    m.key,
			Label:  m.label,
			Status: status,
			Value:  m.stats.Avg,
			Unit:   m.unit,
			Text:   MetricText(m.key, status, high),
		})
	}

	var ralerts []model.ReportAlert
	for _, a := range alerts {
		ralerts = append(ralerts, model.ReportAlert{Level: a.Severity, Message: a.Message})
	}
	return model.ReportDetails{Metrics: metrics, Alerts: ralerts, Note: note}
}
