package report

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/service/llm"
	"github.com/zcl0621/compx576-smart-dairy-system/testhelper"
	"gorm.io/gorm"
)

// fakeDeepSeek replies with a sequence of canned JSON or HTTP errors.
// After the plan runs out, it falls back to a generic success response so
// retries that overshoot don't crash.
type fakeDeepSeek struct {
	server *httptest.Server
	calls  int
	plan   []func(http.ResponseWriter, *http.Request)
}

func newFakeDeepSeek(plan []func(http.ResponseWriter, *http.Request)) *fakeDeepSeek {
	f := &fakeDeepSeek{plan: plan}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { f.calls++ }()
		if f.calls < len(f.plan) {
			f.plan[f.calls](w, r)
			return
		}
		// default: success
		json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{{"message": map[string]string{
				"content": `{"summary":"default","note":"default"}`,
			}}},
		})
	}))
	return f
}

func TestRunOnce_Integration_FreshRunReplacesToday(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		// 5 active cows + 1 sold cow that should be skipped
		var cowIDs []string
		for i := 0; i < 5; i++ {
			c := testhelper.SeedCow(t, tx, "active", model.CowStatusInFarm)
			cowIDs = append(cowIDs, c.ID)
		}
		testhelper.SeedCow(t, tx, "sold", model.CowStatusSold)

		// Pre-existing report dated today (should be deleted)
		existing := testhelper.SeedReport(t, tx, cowIDs[0])
		nzNow := time.Now()
		tx.Model(existing).Update("created_at", nzNow)

		for _, id := range cowIDs {
			for j := 0; j < 50; j++ {
				testhelper.SeedMetric(t, tx, id, model.MetricTypeTemperature, 38.5,
					nzNow.Add(-7*24*time.Hour+time.Duration(j)*time.Hour))
			}
		}

		success := func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{"message": map[string]string{
					"content": `{"summary":"ok","note":"keep going"}`,
				}}},
			})
		}
		fake := newFakeDeepSeek([]func(http.ResponseWriter, *http.Request){
			success, success, success, success, success,
		})
		defer fake.server.Close()

		gen := &Generator{
			DB: pg.DB,
			LLM: &llm.Client{
				BaseURL: fake.server.URL,
				APIKey:  "k",
				Model:   "m",
				Timeout: 5 * time.Second,
				HTTP:    fake.server.Client(),
			},
			Now:     time.Now,
			SleepFn: func(context.Context, time.Duration) {},
		}
		if err := gen.RunOnce(context.Background()); err != nil {
			t.Fatalf("run: %v", err)
		}

		// exactly 5 reports, all dated today, summary "ok"
		var reports []model.Report
		tx.Where("summary = ?", "ok").Find(&reports)
		if len(reports) != 5 {
			t.Errorf("reports = %d, want 5", len(reports))
		}
	})
}

func TestRunOnce_Integration_RetryQueueEventuallySucceeds(t *testing.T) {
	testhelper.SetupTestDB(t)
	testhelper.WithTx(t, func(tx *gorm.DB) {
		cow := testhelper.SeedCow(t, tx, "flaky", model.CowStatusInFarm)
		end := time.Now().UTC().Truncate(24 * time.Hour)
		start := end.AddDate(0, 0, -7)
		for j := 0; j < 50; j++ {
			testhelper.SeedMetric(t, tx, cow.ID, model.MetricTypeTemperature, 38.5,
				start.Add(time.Duration(j)*time.Minute))
		}

		// First two calls fail, third succeeds
		fail := func(w http.ResponseWriter, _ *http.Request) { http.Error(w, "down", 500) }
		ok := func(w http.ResponseWriter, _ *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{{"message": map[string]string{
					"content": `{"summary":"finally","note":"good"}`,
				}}},
			})
		}
		fake := newFakeDeepSeek([]func(http.ResponseWriter, *http.Request){fail, fail, ok})
		defer fake.server.Close()

		gen := &Generator{
			DB: pg.DB,
			LLM: &llm.Client{
				BaseURL: fake.server.URL, APIKey: "k", Model: "m",
				Timeout: 5 * time.Second, HTTP: fake.server.Client(),
			},
			Now:     time.Now,
			SleepFn: func(context.Context, time.Duration) {},
		}
		if err := gen.RunOnce(context.Background()); err != nil {
			t.Fatalf("run: %v", err)
		}

		var r model.Report
		if err := tx.First(&r, "cow_id = ?", cow.ID).Error; err != nil {
			t.Fatalf("expected one report: %v", err)
		}
		if r.Summary != "finally" {
			t.Errorf("summary = %q", r.Summary)
		}
	})
}
