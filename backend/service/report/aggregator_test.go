package report

import (
	"testing"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

func TestAggregateCowMetrics_FullData(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "alpha", model.CowStatusInFarm)
		end := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)
		start := end.AddDate(0, 0, -7)
		baselineEnd := start
		baselineStart := baselineEnd.AddDate(0, 0, -7)

		// current window: 100 temperature readings, all 38.5
		for i := 0; i < 100; i++ {
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeTemperature, 38.5,
				start.Add(time.Duration(i)*time.Minute))
		}
		// baseline window: 50 weight readings of 600.0
		for i := 0; i < 50; i++ {
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeWeight, 600.0,
				baselineStart.Add(time.Duration(i)*time.Hour))
		}
		// current window: 7 weight readings of 612.0 (2% gain)
		for i := 0; i < 7; i++ {
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeWeight, 612.0,
				start.Add(time.Duration(i)*time.Hour*24))
		}

		stats, err := AggregateCowMetrics(pg.DB, cow.ID, start, end, baselineStart, baselineEnd)
		if err != nil {
			t.Fatalf("aggregate: %v", err)
		}
		if stats.Current.Temperature.Count != 100 {
			t.Errorf("current temp count = %d, want 100", stats.Current.Temperature.Count)
		}
		if stats.Current.Temperature.Avg < 38.49 || stats.Current.Temperature.Avg > 38.51 {
			t.Errorf("current temp avg = %f, want ~38.5", stats.Current.Temperature.Avg)
		}
		if stats.Baseline.Weight.Count != 50 {
			t.Errorf("baseline weight count = %d, want 50", stats.Baseline.Weight.Count)
		}
		if stats.Baseline.Weight.Avg < 599.9 || stats.Baseline.Weight.Avg > 600.1 {
			t.Errorf("baseline weight avg = %f, want ~600", stats.Baseline.Weight.Avg)
		}
	})
}

func TestAggregateCowMetrics_NoData(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "ghost", model.CowStatusInFarm)
		end := time.Now()
		start := end.AddDate(0, 0, -7)

		stats, err := AggregateCowMetrics(pg.DB, cow.ID, start, end,
			start.AddDate(0, 0, -7), start)
		if err != nil {
			t.Fatalf("aggregate: %v", err)
		}
		if stats.Current.Temperature.Count != 0 {
			t.Errorf("expected 0 readings, got %d", stats.Current.Temperature.Count)
		}
		if stats.Current.CoveragePct != 0 {
			t.Errorf("expected 0%% coverage, got %f", stats.Current.CoveragePct)
		}
	})
}

func TestAggregateCowMetrics_AbnormalCount(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "fever", model.CowStatusInFarm)
		end := time.Now().UTC()
		start := end.AddDate(0, 0, -7)

		// 10 normal readings, 5 abnormal-high
		for i := 0; i < 10; i++ {
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeTemperature, 38.5,
				start.Add(time.Duration(i)*time.Hour))
		}
		for i := 0; i < 5; i++ {
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeTemperature, 39.5,
				start.Add(time.Duration(i+10)*time.Hour))
		}

		stats, err := AggregateCowMetrics(pg.DB, cow.ID, start, end,
			start.AddDate(0, 0, -7), start)
		if err != nil {
			t.Fatalf("aggregate: %v", err)
		}
		if stats.Current.Temperature.Count != 15 {
			t.Errorf("count = %d, want 15", stats.Current.Temperature.Count)
		}
		if stats.Current.Temperature.AbnormalCount != 5 {
			t.Errorf("abnormal = %d, want 5", stats.Current.Temperature.AbnormalCount)
		}
	})
}

func TestAggregateCowMetrics_Coverage(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "halfdata", model.CowStatusInFarm)
		end := time.Now().UTC()
		start := end.AddDate(0, 0, -7)

		// Expected at 30s intervals: 7 * 24 * 60 * 2 = 20160 per metric.
		// Seed 100 readings per metric → expected coverage ≈ 100/20160 ≈ 0.5%.
		// Math doesn't care about seed count; keep it tiny so the test stays fast.
		for i := 0; i < 100; i++ {
			ts := start.Add(time.Duration(i) * time.Minute)
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeTemperature, 38.5, ts)
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeHeartRate, 70, ts)
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeBloodOxygen, 96, ts)
		}

		stats, err := AggregateCowMetrics(pg.DB, cow.ID, start, end,
			start.AddDate(0, 0, -7), start)
		if err != nil {
			t.Fatalf("aggregate: %v", err)
		}
		if stats.Current.CoveragePct < 0.4 || stats.Current.CoveragePct > 0.6 {
			t.Errorf("coverage = %f, want ~0.5", stats.Current.CoveragePct)
		}
	})
}
