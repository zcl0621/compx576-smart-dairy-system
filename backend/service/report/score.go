package report

import "github.com/zcl0621/compx576-smart-dairy-system/model"

const (
	deductionLight    = 3
	deductionModerate = 8
	deductionSevere   = 15

	deductionAlertWarning  = 3
	deductionAlertCritical = 8
	alertDeductionCap      = 30

	offlineCoverageFloor = 30.0 // % below which the offline penalty kicks in
	offlinePenaltyMul    = 1.5

	minBaselineCount = 3
)

// CalculateScore returns the health score, clamped to [0, 100].
func CalculateScore(stats *CowStats, alerts []model.Alert) float64 {
	deduction := 0.0
	deduction += metricDeviationDeduction(stats)
	deduction += alertDeduction(alerts)
	deduction += offlineDeduction(stats.Current.CoveragePct)

	score := 100.0 - deduction
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func metricDeviationDeduction(stats *CowStats) float64 {
	var d float64
	cur := stats.Current

	// skip metrics with no readings
	if cur.Temperature.Count > 0 {
		d += levelDeduction(ClassifyDeviation(model.MetricTypeTemperature, cur.Temperature.Avg, 0))
	}
	if cur.HeartRate.Count > 0 {
		d += levelDeduction(ClassifyDeviation(model.MetricTypeHeartRate, cur.HeartRate.Avg, 0))
	}
	if cur.BloodOxygen.Count > 0 {
		d += levelDeduction(ClassifyDeviation(model.MetricTypeBloodOxygen, cur.BloodOxygen.Avg, 0))
	}

	// skip weight/milk if baseline is too thin
	if cur.Weight.Count > 0 && stats.Baseline.Weight.Count >= minBaselineCount {
		d += levelDeduction(ClassifyDeviation(model.MetricTypeWeight, cur.Weight.Avg, stats.Baseline.Weight.Avg))
	}
	if cur.MilkAmount.Count > 0 && stats.Baseline.MilkAmount.Count >= minBaselineCount {
		d += levelDeduction(ClassifyDeviation(model.MetricTypeMilkAmount, cur.MilkAmount.Avg, stats.Baseline.MilkAmount.Avg))
	}
	return d
}

func levelDeduction(lv DeviationLevel) float64 {
	switch lv {
	case DeviationLight:
		return deductionLight
	case DeviationModerate:
		return deductionModerate
	case DeviationSevere:
		return deductionSevere
	}
	return 0
}

func alertDeduction(alerts []model.Alert) float64 {
	var d float64
	for _, a := range alerts {
		switch a.Severity {
		case model.AlertSeverityWarning:
			d += deductionAlertWarning
		case model.AlertSeverityCritical:
			d += deductionAlertCritical
		}
	}
	if d > alertDeductionCap {
		return alertDeductionCap
	}
	return d
}

func offlineDeduction(coveragePct float64) float64 {
	offlinePct := 100.0 - coveragePct
	excess := offlinePct - offlineCoverageFloor
	if excess < 0 {
		return 0
	}
	return excess * offlinePenaltyMul
}
