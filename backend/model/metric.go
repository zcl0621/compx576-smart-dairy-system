package model

type Metric struct {
	BaseModel
	CowID       string  `json:"cow_id" gorm:"not null;index"`
	Source      string  `json:"source" gorm:"not null"`      // cow_agent/milking_machine/weight_machine
	MetricType  string  `json:"metric_type" gorm:"not null"` //temperature/heart_rate/blood_oxygen/latitude/longitude/milk_amount/milking_duration/weight
	MetricValue float64 `json:"metric_value" gorm:"type:numeric(10,2);not null"`
	Unit        string  `json:"unit" gorm:"not null"` //celsius/bpm/percent/degrees/liters/seconds/kg
}
