package report

// MetricStats holds aggregate values for one metric over a time window.
// Count == 0 means no readings.
type MetricStats struct {
	Count         int
	Avg           float64
	Min           float64
	Max           float64
	AbnormalCount int
}

// WindowStats holds per-metric stats and health-metric coverage for a window.
type WindowStats struct {
	Temperature MetricStats
	HeartRate   MetricStats
	BloodOxygen MetricStats
	Weight      MetricStats
	MilkAmount  MetricStats
	// CoveragePct is the average coverage across temperature, heart rate, and blood oxygen. Range 0-100.
	CoveragePct float64
}

// CowStats holds the current window and the previous 7-day baseline.
// Score uses Baseline as the weight/milk reference.
type CowStats struct {
	CowID    string
	Current  WindowStats
	Baseline WindowStats
}
