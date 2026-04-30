package report

import (
	"testing"

	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

func TestClassifyDeviation(t *testing.T) {
	cases := []struct {
		name   string
		metric model.MetricType
		avg    float64
		wantLv DeviationLevel
		wantSt model.ReportMetricStatus
	}{
		{"temp normal", model.MetricTypeTemperature, 38.5, DeviationNone, model.ReportMetricStatusNormal},
		{"temp light high", model.MetricTypeTemperature, 39.2, DeviationLight, model.ReportMetricStatusWarning},
		{"temp moderate high", model.MetricTypeTemperature, 39.7, DeviationModerate, model.ReportMetricStatusWarning},
		{"temp severe high", model.MetricTypeTemperature, 40.0, DeviationSevere, model.ReportMetricStatusCritical},
		{"temp light low", model.MetricTypeTemperature, 37.8, DeviationLight, model.ReportMetricStatusWarning},
		{"hr normal", model.MetricTypeHeartRate, 70, DeviationNone, model.ReportMetricStatusNormal},
		{"hr severe", model.MetricTypeHeartRate, 120, DeviationSevere, model.ReportMetricStatusCritical},
		{"spo2 normal", model.MetricTypeBloodOxygen, 95, DeviationNone, model.ReportMetricStatusNormal},
		{"spo2 light", model.MetricTypeBloodOxygen, 89, DeviationLight, model.ReportMetricStatusWarning},
		{"spo2 moderate", model.MetricTypeBloodOxygen, 86, DeviationModerate, model.ReportMetricStatusWarning},
		{"spo2 severe", model.MetricTypeBloodOxygen, 80, DeviationSevere, model.ReportMetricStatusCritical},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lv := ClassifyDeviation(tc.metric, tc.avg, 0)
			if lv != tc.wantLv {
				t.Errorf("level = %v, want %v", lv, tc.wantLv)
			}
			st := DeviationToStatus(lv)
			if st != tc.wantSt {
				t.Errorf("status = %v, want %v", st, tc.wantSt)
			}
		})
	}
}

func TestClassifyDeviation_WeightAgainstBaseline(t *testing.T) {
	// 600 baseline, 615 current = 2.5% gain → none
	if got := ClassifyDeviation(model.MetricTypeWeight, 615, 600); got != DeviationNone {
		t.Errorf("2.5%% gain = %v, want none", got)
	}
	// 600 baseline, 645 current = 7.5% gain → light
	if got := ClassifyDeviation(model.MetricTypeWeight, 645, 600); got != DeviationLight {
		t.Errorf("7.5%% gain = %v, want light", got)
	}
	// 600 baseline, 480 current = 20% loss → moderate (boundary, > 10 ≤ 20)
	if got := ClassifyDeviation(model.MetricTypeWeight, 480, 600); got != DeviationModerate {
		t.Errorf("20%% loss = %v, want moderate", got)
	}
	// 600 baseline, 400 current = 33% loss → severe
	if got := ClassifyDeviation(model.MetricTypeWeight, 400, 600); got != DeviationSevere {
		t.Errorf("33%% loss = %v, want severe", got)
	}
}

func TestMetricText(t *testing.T) {
	if got := MetricText(model.MetricTypeTemperature, model.ReportMetricStatusNormal, false); got != "Normal range" {
		t.Errorf("got %q", got)
	}
	if got := MetricText(model.MetricTypeTemperature, model.ReportMetricStatusWarning, true); got != "Slightly elevated" {
		t.Errorf("got %q", got)
	}
	if got := MetricText(model.MetricTypeTemperature, model.ReportMetricStatusCritical, false); got != "Low" {
		t.Errorf("got %q", got)
	}
	if got := MetricText(model.MetricTypeBloodOxygen, model.ReportMetricStatusCritical, false); got != "Below safe range" {
		t.Errorf("got %q", got)
	}
	if got := MetricText("anything", model.ReportMetricStatusOffline, false); got != "Device offline" {
		t.Errorf("got %q", got)
	}
}
