package report

import (
	"math"

	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

type DeviationLevel int

const (
	DeviationNone DeviationLevel = iota
	DeviationLight
	DeviationModerate
	DeviationSevere
)

// ClassifyDeviation returns how far the metric average deviates from normal.
// For weight and milk, pass the baseline avg; pass 0 for other metrics.
func ClassifyDeviation(metric model.MetricType, avg, baseline float64) DeviationLevel {
	switch metric {
	case model.MetricTypeTemperature:
		return absoluteLevel(avg, 38.0, 39.0, 0.3, 0.8)
	case model.MetricTypeHeartRate:
		return absoluteLevel(avg, 48, 84, 10, 25)
	case model.MetricTypeBloodOxygen:
		// one-sided: only low values are abnormal
		if avg >= 90 {
			return DeviationNone
		}
		if avg >= 88 {
			return DeviationLight
		}
		if avg >= 85 {
			return DeviationModerate
		}
		return DeviationSevere
	case model.MetricTypeWeight:
		return percentLevel(avg, baseline, 5, 10, 20)
	case model.MetricTypeMilkAmount:
		return percentLevel(avg, baseline, 10, 20, 30)
	default:
		return DeviationNone
	}
}

// absoluteLevel classifies distance outside [low, high] into light/moderate/severe.
func absoluteLevel(avg, low, high, light, moderate float64) DeviationLevel {
	var distance float64
	switch {
	case avg < low:
		distance = low - avg
	case avg > high:
		distance = avg - high
	default:
		return DeviationNone
	}
	switch {
	case distance <= light:
		return DeviationLight
	case distance <= moderate:
		return DeviationModerate
	default:
		return DeviationSevere
	}
}

// percentLevel classifies percent deviation from baseline into light/moderate/severe.
func percentLevel(avg, baseline, light, moderate, severe float64) DeviationLevel {
	if baseline == 0 {
		return DeviationNone
	}
	pct := math.Abs(avg-baseline) / baseline * 100
	switch {
	case pct <= light:
		return DeviationNone
	case pct <= moderate:
		return DeviationLight
	case pct <= severe:
		return DeviationModerate
	default:
		return DeviationSevere
	}
}

// DeviationToStatus maps a deviation level to a ReportMetricStatus.
func DeviationToStatus(lv DeviationLevel) model.ReportMetricStatus {
	switch lv {
	case DeviationNone:
		return model.ReportMetricStatusNormal
	case DeviationLight, DeviationModerate:
		return model.ReportMetricStatusWarning
	case DeviationSevere:
		return model.ReportMetricStatusCritical
	}
	return model.ReportMetricStatusNormal
}

// MetricText returns the display label for Metrics[].Text.
// high is true when the average is above normal; ignored for offline/normal status.
func MetricText(metric model.MetricType, status model.ReportMetricStatus, high bool) string {
	if status == model.ReportMetricStatusOffline {
		return "Device offline"
	}
	if status == model.ReportMetricStatusNormal {
		switch metric {
		case model.MetricTypeTemperature:
			return "Normal range"
		case model.MetricTypeHeartRate, model.MetricTypeBloodOxygen:
			return "Normal"
		case model.MetricTypeWeight, model.MetricTypeMilkAmount:
			return "Stable"
		}
		return "Normal"
	}
	// warning and critical use status × direction
	return abnormalText(metric, status, high)
}

type abnormalKey struct {
	metric model.MetricType
	status model.ReportMetricStatus
	high   bool
}

var abnormalTextTable = map[abnormalKey]string{
	{model.MetricTypeTemperature, model.ReportMetricStatusWarning, true}:   "Slightly elevated",
	{model.MetricTypeTemperature, model.ReportMetricStatusWarning, false}:  "Slightly low",
	{model.MetricTypeTemperature, model.ReportMetricStatusCritical, true}:  "Elevated",
	{model.MetricTypeTemperature, model.ReportMetricStatusCritical, false}: "Low",
	{model.MetricTypeHeartRate, model.ReportMetricStatusWarning, true}:     "Slightly elevated",
	{model.MetricTypeHeartRate, model.ReportMetricStatusWarning, false}:    "Slightly low",
	{model.MetricTypeHeartRate, model.ReportMetricStatusCritical, true}:    "Elevated",
	{model.MetricTypeHeartRate, model.ReportMetricStatusCritical, false}:   "Low",
	{model.MetricTypeBloodOxygen, model.ReportMetricStatusWarning, false}:  "Slightly low",
	{model.MetricTypeBloodOxygen, model.ReportMetricStatusCritical, false}: "Below safe range",
	{model.MetricTypeWeight, model.ReportMetricStatusWarning, true}:        "Trending up",
	{model.MetricTypeWeight, model.ReportMetricStatusWarning, false}:       "Trending down",
	{model.MetricTypeWeight, model.ReportMetricStatusCritical, true}:       "Significant gain",
	{model.MetricTypeWeight, model.ReportMetricStatusCritical, false}:      "Significant loss",
	{model.MetricTypeMilkAmount, model.ReportMetricStatusWarning, true}:    "Trending up",
	{model.MetricTypeMilkAmount, model.ReportMetricStatusWarning, false}:   "Trending down",
	{model.MetricTypeMilkAmount, model.ReportMetricStatusCritical, true}:   "Sharp increase",
	{model.MetricTypeMilkAmount, model.ReportMetricStatusCritical, false}:  "Sharp drop",
}

func abnormalText(metric model.MetricType, status model.ReportMetricStatus, high bool) string {
	if v, ok := abnormalTextTable[abnormalKey{metric, status, high}]; ok {
		return v
	}
	return "Abnormal"
}
