package report

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/llm"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

// fakeLLM replays a queue of (output, error) pairs for tests.
type fakeLLM struct {
	answers []llmAnswer
	calls   int
}
type llmAnswer struct {
	out llm.LLMOutput
	err error
}

func (f *fakeLLM) Generate(_ context.Context, _, _ string) (llm.LLMOutput, error) {
	if f.calls >= len(f.answers) {
		return llm.LLMOutput{}, errors.New("no more fake answers")
	}
	a := f.answers[f.calls]
	f.calls++
	return a.out, a.err
}

func TestGenerateOne_InsertsReport(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "alpha", model.CowStatusInFarm)

		// Pin Now and compute window from PeriodWindow so seed timestamps land
		// inside [start, end). Don't reuse a hand-rolled UTC midnight: PeriodWindow
		// is NZ-aligned and the two only coincide for ~12 hours per day.
		fixedNow := time.Date(2026, 4, 30, 5, 0, 0, 0, time.UTC)
		start, _ := PeriodWindow(fixedNow)
		for i := 0; i < 100; i++ {
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeTemperature, 38.5,
				start.Add(time.Duration(i)*time.Minute))
		}

		gen := &Generator{
			DB: pg.DB,
			LLM: &fakeLLM{answers: []llmAnswer{
				{out: llm.LLMOutput{Summary: "ok", Note: "note"}},
			}},
			Now: func() time.Time { return fixedNow },
		}
		if err := gen.GenerateOne(context.Background(), cow.ID); err != nil {
			t.Fatalf("generateOne: %v", err)
		}

		var r model.Report
		if err := tx.Where("cow_id = ?", cow.ID).First(&r).Error; err != nil {
			t.Fatalf("expected one report inserted: %v", err)
		}
		if r.Summary != "ok" {
			t.Errorf("summary = %q", r.Summary)
		}
	})
}

func TestGenerateOne_FailureBubblesUp(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "beta", model.CowStatusInFarm)
		gen := &Generator{
			DB: pg.DB,
			LLM: &fakeLLM{answers: []llmAnswer{
				{err: errors.New("deepseek down")},
			}},
			Now: time.Now,
		}
		if err := gen.GenerateOne(context.Background(), cow.ID); err == nil {
			t.Fatal("expected LLM failure to bubble up")
		}
	})
}

func TestDeleteTodaysReports(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "delta", model.CowStatusInFarm)
		// seed two: one created today (in NZ), one created 3 days ago
		nowNZ := time.Now().In(mustLoadNZ(t))
		old := time.Now().AddDate(0, 0, -3)
		testhelper.SeedReport(t, tx, cow.ID)
		var first model.Report
		tx.Where("cow_id = ?", cow.ID).First(&first)
		tx.Model(&first).Update("created_at", nowNZ)

		second := testhelper.SeedReport(t, tx, cow.ID)
		tx.Model(second).Update("created_at", old)

		if err := DeleteTodaysReports(pg.DB); err != nil {
			t.Fatalf("delete: %v", err)
		}

		var remaining int64
		tx.Model(&model.Report{}).Where("cow_id = ?", cow.ID).Count(&remaining)
		if remaining != 1 {
			t.Errorf("remaining = %d, want 1", remaining)
		}
	})
}

func mustLoadNZ(t *testing.T) *time.Location {
	loc, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		t.Fatalf("load NZ tz: %v", err)
	}
	return loc
}

func TestListActiveCows(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		active := testhelper.SeedCow(t, tx, "alive", model.CowStatusInFarm)
		testhelper.SeedCow(t, tx, "sold", model.CowStatusSold)
		testhelper.SeedCow(t, tx, "inactive", model.CowStatusInactive)

		ids, err := ListActiveCows(pg.DB)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(ids) != 1 || ids[0] != active.ID {
			t.Errorf("ids = %v, want [%s]", ids, active.ID)
		}
	})
}

func TestRunOnce_HappyPath(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		// Pin Now to a fixed instant and seed inside the NZ-aligned window.
		fixedNow := time.Date(2026, 4, 30, 5, 0, 0, 0, time.UTC)
		start, _ := PeriodWindow(fixedNow)
		for i := 0; i < 3; i++ {
			cow := testhelper.SeedCow(t, tx, "cow", model.CowStatusInFarm)
			for j := 0; j < 50; j++ {
				testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeTemperature, 38.5,
					start.Add(time.Duration(j)*time.Minute))
			}
		}

		fake := &fakeLLM{
			answers: []llmAnswer{
				{out: llm.LLMOutput{Summary: "s1", Note: "n1"}},
				{out: llm.LLMOutput{Summary: "s2", Note: "n2"}},
				{out: llm.LLMOutput{Summary: "s3", Note: "n3"}},
			},
		}
		gen := &Generator{DB: pg.DB, LLM: fake, Now: func() time.Time { return fixedNow }}

		err := gen.RunOnce(context.Background())
		if err != nil {
			t.Fatalf("run: %v", err)
		}

		var count int64
		tx.Model(&model.Report{}).Count(&count)
		if count != 3 {
			t.Errorf("reports = %d, want 3", count)
		}
	})
}

func TestRunOnce_FinalFailureCreatesAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "doomed", model.CowStatusInFarm)

		// 11 attempts × 1 cow → 11 failures
		var answers []llmAnswer
		for i := 0; i < 11; i++ {
			answers = append(answers, llmAnswer{err: errors.New("nope")})
		}
		fake := &fakeLLM{answers: answers}
		fixedNow := time.Date(2026, 4, 30, 5, 0, 0, 0, time.UTC)
		gen := &Generator{DB: pg.DB, LLM: fake, Now: func() time.Time { return fixedNow }}
		// cut backoff to zero so the test runs fast
		gen.SleepFn = func(context.Context, time.Duration) {}

		if err := gen.RunOnce(context.Background()); err != nil {
			t.Fatalf("run: %v", err)
		}

		var alerts []model.Alert
		tx.Where("cow_id = ? AND metric_key = ?", cow.ID, model.MetricTypeReportFailure).
			Find(&alerts)
		if len(alerts) != 1 {
			t.Errorf("alerts = %d, want 1", len(alerts))
		}
	})
}

func TestRunOnce_SuccessResolvesPriorFailureAlert(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "recovered", model.CowStatusInFarm)
		// pre-existing active failure alert from a previous day
		testhelper.SeedAlertForMetric(t, tx, cow.ID, model.MetricTypeReportFailure,
			model.AlertSeverityCritical, model.AlertStatusActive)

		fixedNow := time.Date(2026, 4, 30, 5, 0, 0, 0, time.UTC)
		start, _ := PeriodWindow(fixedNow)
		for j := 0; j < 50; j++ {
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeTemperature, 38.5,
				start.Add(time.Duration(j)*time.Minute))
		}

		gen := &Generator{
			DB:  pg.DB,
			LLM: &fakeLLM{answers: []llmAnswer{{out: llm.LLMOutput{Summary: "ok", Note: "n"}}}},
			Now: func() time.Time { return fixedNow },
		}
		if err := gen.RunOnce(context.Background()); err != nil {
			t.Fatalf("run: %v", err)
		}

		var a model.Alert
		if err := tx.Where("cow_id = ? AND metric_key = ?", cow.ID, model.MetricTypeReportFailure).
			First(&a).Error; err != nil {
			t.Fatalf("load alert: %v", err)
		}
		if a.Status != model.AlertStatusResolved {
			t.Errorf("alert status = %s, want resolved", a.Status)
		}
	})
}
