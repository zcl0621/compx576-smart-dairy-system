package report

import (
	"testing"

	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

func makeStats(curTempAvg, curTempCount float64) *CowStats {
	return &CowStats{
		Current: WindowStats{
			Temperature: MetricStats{Count: int(curTempCount), Avg: curTempAvg},
			HeartRate:   MetricStats{Count: int(curTempCount), Avg: 70},
			BloodOxygen: MetricStats{Count: int(curTempCount), Avg: 96},
			CoveragePct: 100,
		},
	}
}

func TestCalculateScore_AllNormal(t *testing.T) {
	stats := makeStats(38.5, 1000)
	s := CalculateScore(stats, nil)
	if s != 100 {
		t.Errorf("score = %f, want 100", s)
	}
}

func TestCalculateScore_LightTempDeviation(t *testing.T) {
	stats := makeStats(39.2, 1000) // 0.2 outside → light → -3
	s := CalculateScore(stats, nil)
	if s != 97 {
		t.Errorf("score = %f, want 97", s)
	}
}

func TestCalculateScore_SevereTempDeviation(t *testing.T) {
	stats := makeStats(40.0, 1000) // 1.0 outside → severe → -15
	s := CalculateScore(stats, nil)
	if s != 85 {
		t.Errorf("score = %f, want 85", s)
	}
}

func TestCalculateScore_AlertDeductionWithCap(t *testing.T) {
	stats := makeStats(38.5, 1000)
	var alerts []model.Alert
	// 50 warnings would be -150 uncapped → cap at -30
	for i := 0; i < 50; i++ {
		alerts = append(alerts, model.Alert{Severity: model.AlertSeverityWarning})
	}
	s := CalculateScore(stats, alerts)
	if s != 70 {
		t.Errorf("score = %f, want 70 (cap at -30)", s)
	}
}

func TestCalculateScore_OfflinePenalty(t *testing.T) {
	stats := &CowStats{
		Current: WindowStats{
			Temperature: MetricStats{Count: 0},
			HeartRate:   MetricStats{Count: 0},
			BloodOxygen: MetricStats{Count: 0},
			CoveragePct: 50, // 50% offline → (50-30)*1.5 = 30 deduction
		},
	}
	s := CalculateScore(stats, nil)
	if s != 70 {
		t.Errorf("score = %f, want 70", s)
	}
}

func TestCalculateScore_FullyOffline(t *testing.T) {
	stats := &CowStats{
		Current: WindowStats{
			Temperature: MetricStats{Count: 0},
			HeartRate:   MetricStats{Count: 0},
			BloodOxygen: MetricStats{Count: 0},
			CoveragePct: 0,
		},
	}
	s := CalculateScore(stats, nil)
	if s != 0 {
		t.Errorf("score = %f, want 0", s)
	}
}

func TestCalculateScore_FloorAtZero(t *testing.T) {
	// Drive total deduction past 100 to verify the floor clamp.
	// - severe temp deviation: -15
	// - 5 critical alerts (cap at -30): -30
	// - 95% offline → (95-30)*1.5: -97.5
	// Total deduction 142.5 → 100 - 142.5 = -42.5 → clamp to 0.
	stats := makeStats(42.0, 1000)
	stats.Current.CoveragePct = 5
	var alerts []model.Alert
	for i := 0; i < 5; i++ {
		alerts = append(alerts, model.Alert{Severity: model.AlertSeverityCritical})
	}
	s := CalculateScore(stats, alerts)
	if s != 0 {
		t.Errorf("score = %f, want 0 (clamped from -42.5)", s)
	}
}

func TestCalculateScore_MissingBaselineSkipsWeight(t *testing.T) {
	stats := &CowStats{
		Current: WindowStats{
			Temperature: MetricStats{Count: 100, Avg: 38.5},
			HeartRate:   MetricStats{Count: 100, Avg: 70},
			BloodOxygen: MetricStats{Count: 100, Avg: 96},
			Weight:      MetricStats{Count: 7, Avg: 1000.0}, // current shows huge weight
			CoveragePct: 100,
		},
		Baseline: WindowStats{
			Weight: MetricStats{Count: 1, Avg: 600.0}, // < 3 readings → skip deduction
		},
	}
	s := CalculateScore(stats, nil)
	if s != 100 {
		t.Errorf("score = %f, want 100 (weight skipped due to thin baseline)", s)
	}
}
