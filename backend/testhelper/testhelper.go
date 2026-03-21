// Package testhelper holds shared test setup and seed helpers
// Each test opens tx on real db and rolls it back when done
// so data does not leak between runs
package testhelper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	redisdb "github.com/zcl0621/compx576-smart-dairy-system/db/redis"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/util"
	"gorm.io/gorm"
)

var setupOnce sync.Once

// init config, logger and pg db once per test binary
// change workdir to backend root so config.yaml is found
func SetupTestDB(t *testing.T) {
	t.Helper()
	setupOnce.Do(func() {
		// one level up is backend
		_, thisFile, _, _ := runtime.Caller(0)
		backendDir := filepath.Dir(filepath.Dir(thisFile))
		if err := os.Chdir(backendDir); err != nil {
			t.Fatalf("chdir to backend: %v", err)
		}

		if err := config.InitConfig(); err != nil {
			t.Fatalf("init config: %v", err)
		}
		if err := projectlog.Init(); err != nil {
			t.Fatalf("init log: %v", err)
		}
		if err := pg.InitDB(); err != nil {
			t.Fatalf("init db: %v", err)
		}
		if err := redisdb.InitRedis(); err != nil {
			t.Fatalf("init redis: %v", err)
		}
		// clean all tables so tests start empty
		if err := pg.DB.Exec("TRUNCATE TABLE users, cows, reports, metrics, alerts CASCADE").Error; err != nil {
			t.Fatalf("truncate tables: %v", err)
		}
		flushRedis(t)
	})
}

// run fn in tx and always roll it back
func WithTx(t *testing.T, fn func(tx *gorm.DB)) {
	t.Helper()
	flushRedis(t)

	tx := pg.DB.Begin()
	if tx.Error != nil {
		t.Fatalf("begin tx: %v", tx.Error)
	}

	original := pg.DB
	pg.DB = tx

	defer func() {
		pg.DB = original
		tx.Rollback() //nolint:errcheck
		flushRedis(t)
	}()

	fn(tx)
}

func FlushRedis(t *testing.T) {
	t.Helper()
	flushRedis(t)
}

func flushRedis(t *testing.T) {
	t.Helper()
	if redisdb.GetClient() == nil {
		return
	}
	if err := redisdb.FlushDB(); err != nil {
		t.Fatalf("flush redis: %v", err)
	}
}

// seed user and return model
func SeedUser(t *testing.T, db *gorm.DB, username, email, password string) *model.User {
	t.Helper()
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	u := &model.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
	}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u
}

// seed cow and return model
func SeedCow(t *testing.T, db *gorm.DB, name string, status model.CowStatus, condition model.CowCondition) *model.Cow {
	t.Helper()
	tag := fmt.Sprintf("TAG-%d", time.Now().UnixNano())
	c := &model.Cow{
		Name:      name,
		Tag:       tag,
		Age:       2,
		Status:    status,
		Condition: condition,
	}
	if err := db.Create(c).Error; err != nil {
		t.Fatalf("seed cow: %v", err)
	}
	return c
}

// seed metric row with created_at
func SeedMetric(t *testing.T, db *gorm.DB, cowID string, metricType model.MetricType, value float64, createdAt time.Time) *model.Metric {
	t.Helper()
	m := &model.Metric{
		CowID:       cowID,
		Source:      model.MetricSourceCowAgent,
		MetricType:  metricType,
		MetricValue: value,
		Unit:        unitFor(metricType),
	}
	if err := db.Create(m).Error; err != nil {
		t.Fatalf("seed metric: %v", err)
	}
	// set created_at so range query works
	if err := db.Model(m).Update("created_at", createdAt).Error; err != nil {
		t.Fatalf("seed metric set created_at: %v", err)
	}
	m.CreatedAt = createdAt
	return m
}

// seed alert and return model
func SeedAlert(t *testing.T, db *gorm.DB, cowID string, severity model.AlertSeverity, status model.AlertStatus) *model.Alert {
	t.Helper()
	a := &model.Alert{
		CowID:     cowID,
		MetricKey: model.MetricTypeTemperature,
		Title:     "test alert",
		Message:   fmt.Sprintf("%s alert for %s", severity, cowID),
		Severity:  severity,
		Status:    status,
	}
	if err := db.Create(a).Error; err != nil {
		t.Fatalf("seed alert: %v", err)
	}
	return a
}

// seed report and return model
func SeedReport(t *testing.T, db *gorm.DB, cowID string) *model.Report {
	t.Helper()
	now := time.Now()
	r := &model.Report{
		CowID:       cowID,
		PeriodStart: now.AddDate(0, 0, -7),
		PeriodEnd:   now,
		Summary:     "test summary",
		Score:       80,
		Details:     model.ReportDetails{Note: "test"},
	}
	if err := db.Create(r).Error; err != nil {
		t.Fatalf("seed report: %v", err)
	}
	return r
}

func unitFor(t model.MetricType) model.MetricUnit {
	switch t {
	case model.MetricTypeTemperature:
		return model.MetricUnitCelsius
	case model.MetricTypeHeartRate:
		return model.MetricUnitBPM
	case model.MetricTypeBloodOxygen:
		return model.MetricUnitPercent
	case model.MetricTypeLatitude, model.MetricTypeLongitude:
		return model.MetricUnitDegrees
	case model.MetricTypeMilkAmount:
		return model.MetricUnitLiters
	case model.MetricTypeWeight:
		return model.MetricUnitKG
	default:
		return model.MetricUnitCelsius
	}
}
