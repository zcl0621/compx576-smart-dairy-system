package model

type MetricSource string

const (
	MetricSourceCowAgent       MetricSource = "cow_agent"
	MetricSourceMilkingMachine MetricSource = "milking_machine"
	MetricSourceWeightMachine  MetricSource = "weight_machine"
)

type MetricType string

type MetricRange string

const (
	MetricTypeTemperature     MetricType = "temperature"
	MetricTypeHeartRate       MetricType = "heart_rate"
	MetricTypeBloodOxygen     MetricType = "blood_oxygen"
	MetricTypeLatitude        MetricType = "latitude"
	MetricTypeLongitude       MetricType = "longitude"
	MetricTypeMilkAmount      MetricType = "milk_amount"
	MetricTypeMilkingDuration MetricType = "milking_duration"
	MetricTypeWeight          MetricType = "weight"
	MetricTypeDevice          MetricType = "device"
	MetricTypeReportFailure   MetricType = "report_failure"
)

const (
	MetricRange24H MetricRange = "24h"
	MetricRange7D  MetricRange = "7d"
	MetricRange30D MetricRange = "30d"
	MetricRangeAll MetricRange = "all"
)

type MetricUnit string

const (
	MetricUnitCelsius MetricUnit = "celsius"
	MetricUnitBPM     MetricUnit = "bpm"
	MetricUnitPercent MetricUnit = "percent"
	MetricUnitDegrees MetricUnit = "degrees"
	MetricUnitLiters  MetricUnit = "liters"
	MetricUnitSeconds MetricUnit = "seconds"
	MetricUnitKG      MetricUnit = "kg"
)

type Metric struct {
	BaseModel
	CowID       string       `json:"cow_id" gorm:"not null;index"`
	Source      MetricSource `json:"source" gorm:"not null"`
	MetricType  MetricType   `json:"metric_type" gorm:"not null"`
	MetricValue float64      `json:"metric_value" gorm:"type:numeric(10,2);not null"`
	Unit        MetricUnit   `json:"unit" gorm:"not null"`
}
