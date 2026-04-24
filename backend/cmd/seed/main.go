// seed inserts mock data into the database for development
// run: go run cmd/seed/main.go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/config"
	"github.com/zcl0621/compx576-smart-dairy-system/db/pg"
	projectlog "github.com/zcl0621/compx576-smart-dairy-system/log"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
	"github.com/zcl0621/compx576-smart-dairy-system/util"
	"gorm.io/gorm"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	if err := projectlog.Init(); err != nil {
		panic(err)
	}
	if err := pg.InitDB(); err != nil {
		panic(err)
	}
	if err := model.Migrate(pg.DB); err != nil {
		panic(err)
	}

	// clean existing data
	pg.DB.Exec("TRUNCATE TABLE users, cows, reports, metrics, alerts CASCADE")
	fmt.Println("tables truncated")

	seedUsers(pg.DB)
	cowIDs := seedCows(pg.DB)
	seedMetrics(pg.DB, cowIDs)
	seedAlerts(pg.DB, cowIDs)
	seedReports(pg.DB, cowIDs)

	fmt.Println("seed done")
}

func seedUsers(db *gorm.DB) {
	users := []struct {
		username string
		email    string
		password string
	}{
		{"admin", "admin@smartdairy.local", "password123"},
		{"farmworker", "worker@smartdairy.local", "password123"},
		{"zcl0621", "zcl0621@hotmail.com", "zxc123456"},
	}

	for _, u := range users {
		hashed, err := util.HashPassword(u.password)
		if err != nil {
			panic(err)
		}
		if err := db.Create(&model.User{
			Username: u.username,
			Password: hashed,
			Email:    u.email,
		}).Error; err != nil {
			panic(err)
		}
	}
	fmt.Println("seeded 3 users")
}

func seedCows(db *gorm.DB) []string {
	type cowDef struct {
		name       string
		tag        string
		age        int
		canMilking bool
		status     model.CowStatus
	}

	defs := []cowDef{
		{"Bessie", "DC-2401", 5, true, model.CowStatusInFarm},
		{"Daisy", "DC-2402", 4, true, model.CowStatusInFarm},
		{"Clara", "DC-2403", 6, false, model.CowStatusInFarm},
		{"Molly", "DC-2404", 4, true, model.CowStatusInFarm},
		{"Ella", "DC-2405", 3, true, model.CowStatusInFarm},
		{"Lucy", "DC-2406", 5, true, model.CowStatusInFarm},
		{"Bella", "DC-2407", 4, true, model.CowStatusInFarm},
		{"Rose", "DC-2408", 6, true, model.CowStatusInFarm},
		{"Sophie", "DC-2409", 3, true, model.CowStatusInFarm},
		{"Rosie", "DC-2410", 5, true, model.CowStatusInFarm},
		{"Maggie", "DC-2411", 4, true, model.CowStatusInFarm},
		{"Buttercup", "DC-2412", 7, false, model.CowStatusInFarm},
		{"Penny", "DC-2413", 3, true, model.CowStatusInFarm},
		{"Hazel", "DC-2414", 5, true, model.CowStatusInFarm},
		{"Poppy", "DC-2415", 4, true, model.CowStatusInFarm},
		{"Clover", "DC-2416", 5, true, model.CowStatusInFarm},
		{"Willow", "DC-2417", 6, false, model.CowStatusInFarm},
		{"Ginger", "DC-2418", 4, true, model.CowStatusInFarm},
		{"Honey", "DC-2419", 3, true, model.CowStatusInFarm},
		{"Luna", "DC-2420", 5, true, model.CowStatusInFarm},
		{"Annie", "DC-2421", 4, true, model.CowStatusInFarm},
		{"Ruby", "DC-2422", 6, true, model.CowStatusInFarm},
		{"Pearl", "DC-2423", 5, true, model.CowStatusInFarm},
		{"Nora", "DC-2424", 4, true, model.CowStatusInFarm},
		{"Olive", "DC-2425", 3, true, model.CowStatusInFarm},
	}

	ids := make([]string, 0, len(defs))
	for _, d := range defs {
		cow := &model.Cow{
			Name:       d.name,
			Tag:        d.tag,
			Age:        d.age,
			CanMilking: d.canMilking,
			Status:     d.status,
		}
		if err := db.Create(cow).Error; err != nil {
			panic(err)
		}
		ids = append(ids, cow.ID)
	}
	fmt.Printf("seeded %d cows\n", len(ids))
	return ids
}

func seedMetrics(db *gorm.DB, cowIDs []string) {
	now := time.Now()
	count := 0

	// generate 7 days of data, one record every 5 minutes
	for _, cowID := range cowIDs {
		var batch []model.Metric
		for minutes := 7 * 24 * 60; minutes >= 0; minutes -= 5 {
			ts := now.Add(-time.Duration(minutes) * time.Minute)
			hour := ts.Hour()

			// temperature: base 38.5, slight daily cycle
			temp := 38.5 + 0.3*math.Sin(float64(hour)/24*2*math.Pi) + randFloat(-0.2, 0.2)
			batch = append(batch, metricRow(cowID, model.MetricTypeTemperature, model.MetricSourceCowAgent, round2(temp), model.MetricUnitCelsius, ts))

			// heart rate: base 72
			hr := 72.0 + 5*math.Sin(float64(hour)/24*2*math.Pi) + randFloat(-3, 3)
			batch = append(batch, metricRow(cowID, model.MetricTypeHeartRate, model.MetricSourceCowAgent, round2(hr), model.MetricUnitBPM, ts))

			// blood oxygen: base 96
			bo := 96.0 + randFloat(-1.5, 1.5)
			batch = append(batch, metricRow(cowID, model.MetricTypeBloodOxygen, model.MetricSourceCowAgent, round2(bo), model.MetricUnitPercent, ts))

			// latitude/longitude: small drift around Waikato area
			lat := -37.7833 + randFloat(-0.002, 0.002)
			lon := 175.2833 + randFloat(-0.002, 0.002)
			batch = append(batch, metricRow(cowID, model.MetricTypeLatitude, model.MetricSourceCowAgent, round2(lat), model.MetricUnitDegrees, ts))
			batch = append(batch, metricRow(cowID, model.MetricTypeLongitude, model.MetricSourceCowAgent, round2(lon), model.MetricUnitDegrees, ts))

			// weight: once per day (at day boundary of the 5-min loop)
			if minutes%(24*60) == 0 {
				w := 600.0 + randFloat(-15, 15)
				batch = append(batch, metricRow(cowID, model.MetricTypeWeight, model.MetricSourceWeightMachine, round2(w), model.MetricUnitKG, ts))
			}

			// milk: twice per day (at day boundary and half-day boundary)
			if minutes%(24*60) == 0 || minutes%(24*60) == 12*60 {
				milk := 12.0 + randFloat(-2, 2)
				batch = append(batch, metricRow(cowID, model.MetricTypeMilkAmount, model.MetricSourceMilkingMachine, round2(milk), model.MetricUnitLiters, ts))
			}
		}

		// batch insert in chunks of 2000
		for i := 0; i < len(batch); i += 2000 {
			end := min(i+2000, len(batch))
			if err := db.Create(batch[i:end]).Error; err != nil {
				panic(err)
			}
		}
		count += len(batch)
	}
	fmt.Printf("seeded %d metric rows\n", count)
}

func seedAlerts(db *gorm.DB, cowIDs []string) {
	now := time.Now()
	resolved := now.Add(-3 * time.Hour)

	alertDefs := []struct {
		cowIdx    int
		metricKey model.MetricType
		title     string
		message   string
		severity  model.AlertSeverity
		status    model.AlertStatus
		resolved  *time.Time
	}{
		{0, model.MetricTypeTemperature, "High temperature", "Temperature stayed above normal for 35 minutes", model.AlertSeverityWarning, model.AlertStatusActive, nil},
		{2, model.MetricTypeBloodOxygen, "Low blood oxygen", "Blood oxygen dropped below 90 percent", model.AlertSeverityCritical, model.AlertStatusActive, nil},
		{4, model.MetricTypeDevice, "Device offline", "No signal from collar device for 4 hours", model.AlertSeverityOffline, model.AlertStatusActive, nil},
		{3, model.MetricTypeTemperature, "Temperature back to normal", "Temperature returned to normal range", model.AlertSeverityWarning, model.AlertStatusResolved, &resolved},
		{5, model.MetricTypeHeartRate, "High heart rate", "Heart rate above 100 bpm for over 20 minutes", model.AlertSeverityWarning, model.AlertStatusActive, nil},
		{22, model.MetricTypeBloodOxygen, "Low blood oxygen", "Blood oxygen dropped to 88 percent", model.AlertSeverityCritical, model.AlertStatusActive, nil},
		{16, model.MetricTypeDevice, "Device offline", "No signal from collar device for 5 hours", model.AlertSeverityOffline, model.AlertStatusActive, nil},
		{7, model.MetricTypeTemperature, "High temperature", "Temperature at 39.8 degrees for 40 minutes", model.AlertSeverityWarning, model.AlertStatusResolved, &resolved},
		{6, model.MetricTypeHeartRate, "Irregular heart rate", "Heart rate pattern abnormal in last 30 minutes", model.AlertSeverityCritical, model.AlertStatusActive, nil},
		{15, model.MetricTypeDevice, "Device battery low", "Collar battery below 10 percent", model.AlertSeverityWarning, model.AlertStatusActive, nil},
		{10, model.MetricTypeBloodOxygen, "Blood oxygen back to normal", "Blood oxygen recovered to 96 percent", model.AlertSeverityWarning, model.AlertStatusResolved, &resolved},
	}

	for _, d := range alertDefs {
		a := &model.Alert{
			CowID:      cowIDs[d.cowIdx],
			MetricKey:  d.metricKey,
			Title:      d.title,
			Message:    d.message,
			Severity:   d.severity,
			Status:     d.status,
			ResolvedAt: d.resolved,
		}
		if err := db.Create(a).Error; err != nil {
			panic(err)
		}
	}
	fmt.Printf("seeded %d alerts\n", len(alertDefs))
}

func seedReports(db *gorm.DB, cowIDs []string) {
	now := time.Now()
	start := now.AddDate(0, 0, -7)

	reportDefs := []struct {
		cowIdx  int
		summary string
		score   float64
		details model.ReportDetails
	}{
		{3, "Condition stays stable. Milk amount is good and daily movement is normal.", 88, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusNormal, Value: 38.5, Unit: model.MetricUnitCelsius, Text: "Normal range"},
				{Key: model.MetricTypeHeartRate, Label: "Heart Rate", Status: model.ReportMetricStatusNormal, Value: 72, Unit: model.MetricUnitBPM, Text: "Stable"},
			},
			Alerts: nil,
			Note:   "No unusual temperature trend. Keep current feed plan.",
		}},
		{13, "Need closer watch on temperature and water intake.", 73, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusWarning, Value: 39.1, Unit: model.MetricUnitCelsius, Text: "Slightly elevated"},
				{Key: model.MetricTypeHeartRate, Label: "Heart Rate", Status: model.ReportMetricStatusNormal, Value: 77, Unit: model.MetricUnitBPM, Text: "Normal"},
			},
			Alerts: []model.ReportAlert{{Level: model.AlertSeverityWarning, Message: "Temperature elevated"}},
			Note:   "Check device and body temperature again.",
		}},
		{2, "Critical pattern detected. Need farm staff follow-up soon.", 59, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeBloodOxygen, Label: "Blood Oxygen", Status: model.ReportMetricStatusCritical, Value: 88, Unit: model.MetricUnitPercent, Text: "Below safe range"},
				{Key: model.MetricTypeHeartRate, Label: "Heart Rate", Status: model.ReportMetricStatusWarning, Value: 88, Unit: model.MetricUnitBPM, Text: "Elevated"},
			},
			Alerts: []model.ReportAlert{
				{Level: model.AlertSeverityCritical, Message: "Blood oxygen dropped below 90"},
				{Level: model.AlertSeverityWarning, Message: "Heart rate elevated"},
			},
			Note: "Do physical check.",
		}},
		{6, "Slightly elevated temperature. Monitor for next 24 hours.", 72, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusWarning, Value: 39.0, Unit: model.MetricUnitCelsius, Text: "Trending upward"},
			},
			Note: "Temperature trending upward. Heart rate normal. Blood oxygen stable.",
		}},
		{7, "Critical blood oxygen levels. Immediate attention required.", 65, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeBloodOxygen, Label: "Blood Oxygen", Status: model.ReportMetricStatusCritical, Value: 87, Unit: model.MetricUnitPercent, Text: "Stayed low"},
			},
			Note: "Need urgent physical check and device verification.",
		}},
		{8, "Excellent health. All metrics within normal range.", 92, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusNormal, Value: 38.5, Unit: model.MetricUnitCelsius, Text: "Normal"},
				{Key: model.MetricTypeHeartRate, Label: "Heart Rate", Status: model.ReportMetricStatusNormal, Value: 68, Unit: model.MetricUnitBPM, Text: "Normal"},
				{Key: model.MetricTypeBloodOxygen, Label: "Blood Oxygen", Status: model.ReportMetricStatusNormal, Value: 97, Unit: model.MetricUnitPercent, Text: "Normal"},
			},
			Note: "Stable readings. No action needed.",
		}},
		{9, "Good health with stable vitals and strong production.", 88, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusNormal, Value: 38.6, Unit: model.MetricUnitCelsius, Text: "Normal"},
			},
			Note: "Continue current feeding and observation plan.",
		}},
		{10, "Elevated heart rate detected. Continue monitoring.", 78, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeHeartRate, Label: "Heart Rate", Status: model.ReportMetricStatusWarning, Value: 82, Unit: model.MetricUnitBPM, Text: "Elevated"},
			},
			Note: "Need closer watch on heart rate pattern during next 24 hours.",
		}},
		{11, "Temperature slightly elevated. Continue observation.", 75, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusWarning, Value: 38.9, Unit: model.MetricUnitCelsius, Text: "Slightly high"},
			},
			Note: "Early warning only. Keep routine checks.",
		}},
		{12, "Excellent overall health with strong production.", 90, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusNormal, Value: 38.6, Unit: model.MetricUnitCelsius, Text: "Normal"},
			},
			Note: "No anomaly. Good milk output and stable movement.",
		}},
		{14, "Good condition. All vitals in range.", 85, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusNormal, Value: 38.5, Unit: model.MetricUnitCelsius, Text: "Normal"},
			},
			Note: "Keep current plan.",
		}},
		{0, "Warning pattern in temperature. Watch closely.", 70, model.ReportDetails{
			Metrics: []model.ReportMetric{
				{Key: model.MetricTypeTemperature, Label: "Temperature", Status: model.ReportMetricStatusWarning, Value: 39.2, Unit: model.MetricUnitCelsius, Text: "Elevated"},
			},
			Alerts: []model.ReportAlert{{Level: model.AlertSeverityWarning, Message: "Temperature elevated to 39.2"}},
			Note:   "Temperature stayed above normal. Check cow condition.",
		}},
	}

	for _, d := range reportDefs {
		r := &model.Report{
			CowID:       cowIDs[d.cowIdx],
			PeriodStart: start,
			PeriodEnd:   now,
			Summary:     d.summary,
			Score:       d.score,
			Details:     d.details,
		}
		if err := db.Create(r).Error; err != nil {
			panic(err)
		}
	}
	fmt.Printf("seeded %d reports\n", len(reportDefs))
}

func metricRow(cowID string, mt model.MetricType, src model.MetricSource, value float64, unit model.MetricUnit, ts time.Time) model.Metric {
	return model.Metric{
		BaseModel:   model.BaseModel{CreatedAt: ts, UpdatedAt: ts},
		CowID:       cowID,
		Source:      src,
		MetricType:  mt,
		MetricValue: value,
		Unit:        unit,
	}
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
