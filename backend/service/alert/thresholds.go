package alert

// Temperature thresholds (°C)
const (
	TempCritHigh = 40.5
	TempCritLow  = 37.0
	TempWarnHigh = 39.5
	TempWarnLow  = 38.0
)

// Heart rate thresholds (bpm)
const (
	HRCritHigh = 100.0
	HRCritLow  = 38.0
	HRWarnHigh = 84.0
	HRWarnLow  = 48.0
)

// Blood oxygen thresholds (%)
const (
	BOCrit = 88.0
	BOWarn = 95.0
)

// Milking thresholds (fraction of 7-day average)
const (
	MilkWarnFraction    = 0.70 // below this → warning
	MilkResolveFraction = 0.85 // at or above this → resolve (hysteresis)
)

// Weight deviation thresholds (fraction of 7-day average)
const (
	WeightWarnDeviation    = 0.30 // abs deviation above this → warning
	WeightResolveDeviation = 0.15 // abs deviation below this → resolve (hysteresis)
)

const OfflineThreshold = 10 // minutes without a cow_agent metric
