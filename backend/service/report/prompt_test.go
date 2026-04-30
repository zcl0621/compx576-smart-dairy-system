package report

import (
	"strings"
	"testing"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

func TestBuildSystemPrompt_StableAcrossCalls(t *testing.T) {
	a := BuildSystemPrompt()
	b := BuildSystemPrompt()
	if a != b {
		t.Fatal("system prompt must be identical across calls (caching depends on it)")
	}
	if !strings.Contains(strings.ToLower(a), "json") {
		t.Error("system prompt must mention json (DeepSeek requirement)")
	}
}

func TestBuildUserPrompt_FullData(t *testing.T) {
	cow := &model.Cow{Name: "Bessie"}
	stats := &CowStats{
		Current: WindowStats{
			Temperature: MetricStats{Count: 1000, Avg: 38.5, Min: 38.0, Max: 39.0, AbnormalCount: 0},
			HeartRate:   MetricStats{Count: 1000, Avg: 70, Min: 60, Max: 80, AbnormalCount: 0},
			BloodOxygen: MetricStats{Count: 1000, Avg: 96, Min: 92, Max: 99, AbnormalCount: 0},
			Weight:      MetricStats{Count: 7, Avg: 615, Min: 612, Max: 618},
			MilkAmount:  MetricStats{Count: 14, Avg: 25, Min: 22, Max: 28},
			CoveragePct: 95,
		},
		Baseline: WindowStats{
			Weight:     MetricStats{Count: 7, Avg: 600},
			MilkAmount: MetricStats{Count: 14, Avg: 24},
		},
	}
	var alerts []model.Alert
	start := time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)

	p := BuildUserPrompt(cow, stats, alerts, 92, start, end)
	if !strings.Contains(p, "Bessie") {
		t.Error("user prompt missing cow name")
	}
	if !strings.Contains(p, "Health score: 92") {
		t.Error("user prompt missing score")
	}
	if !strings.Contains(p, "2026-04-23") || !strings.Contains(p, "2026-04-30") {
		t.Error("user prompt missing period dates")
	}
	if !strings.Contains(p, "Alerts in period") {
		t.Error("user prompt missing alerts header")
	}
	if !strings.Contains(p, "+15.0 kg") && !strings.Contains(p, "+15") {
		t.Error("user prompt missing weight delta")
	}
}

func TestBuildUserPrompt_NoBaselineForWeight(t *testing.T) {
	cow := &model.Cow{Name: "NewCow"}
	stats := &CowStats{
		Current: WindowStats{
			Temperature: MetricStats{Count: 1000, Avg: 38.5},
			HeartRate:   MetricStats{Count: 1000, Avg: 70},
			BloodOxygen: MetricStats{Count: 1000, Avg: 96},
			Weight:      MetricStats{Count: 7, Avg: 600},
			CoveragePct: 100,
		},
		Baseline: WindowStats{
			Weight: MetricStats{Count: 0},
		},
	}
	p := BuildUserPrompt(cow, stats, nil, 100,
		time.Now().AddDate(0, 0, -7), time.Now())
	if !strings.Contains(p, "no baseline available") {
		t.Errorf("expected weight to mention 'no baseline available', got:\n%s", p)
	}
}

func TestBuildUserPrompt_NoData(t *testing.T) {
	cow := &model.Cow{Name: "Ghost"}
	stats := &CowStats{Current: WindowStats{CoveragePct: 0}}
	p := BuildUserPrompt(cow, stats, nil, 0,
		time.Now().AddDate(0, 0, -7), time.Now())
	if !strings.Contains(p, "Ghost") {
		t.Error("missing cow name")
	}
	if !strings.Contains(p, "No metric data") {
		t.Errorf("expected no-data message, got:\n%s", p)
	}
	if !strings.Contains(strings.ToLower(p), "json") {
		t.Error("prompt must mention json")
	}
}

func TestBuildUserPrompt_PrintsAllAlertsReceived(t *testing.T) {
	// The caller (GenerateOne) already caps alerts at 20 via .Limit(20).
	// The prompt prints every alert it receives without re-capping.
	cow := &model.Cow{Name: "Loud"}
	stats := &CowStats{
		Current: WindowStats{
			Temperature: MetricStats{Count: 100, Avg: 38.5},
			HeartRate:   MetricStats{Count: 100, Avg: 70},
			BloodOxygen: MetricStats{Count: 100, Avg: 96},
			CoveragePct: 100,
		},
	}
	var alerts []model.Alert
	base := time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC)
	for i := range 20 {
		a := model.Alert{Severity: model.AlertSeverityWarning, Message: "spike"}
		a.CreatedAt = base.Add(time.Duration(i) * time.Hour)
		alerts = append(alerts, a)
	}
	p := BuildUserPrompt(cow, stats, alerts, 80, base, base.AddDate(0, 0, 7))
	count := strings.Count(p, "spike")
	if count != 20 {
		t.Errorf("expected 20 alert lines, got %d:\n%s", count, p)
	}
}
