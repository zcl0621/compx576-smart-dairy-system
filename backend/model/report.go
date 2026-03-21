package model

import (
	"time"
)

type ReportMetricStatus string

const (
	ReportMetricStatusNormal   ReportMetricStatus = "normal"
	ReportMetricStatusWarning  ReportMetricStatus = "warning"
	ReportMetricStatusCritical ReportMetricStatus = "critical"
	ReportMetricStatusOffline  ReportMetricStatus = "offline"
)

type ReportMetric struct {
	Key    MetricType         `json:"key"`
	Label  string             `json:"label"`
	Status ReportMetricStatus `json:"status"`
	Value  float64            `json:"value"`
	Unit   MetricUnit         `json:"unit"`
	Text   string             `json:"text"`
}

type ReportAlert struct {
	Level   AlertSeverity `json:"level"`
	Message string        `json:"message"`
}

type ReportDetails struct {
	Metrics []ReportMetric `json:"metrics"`
	Alerts  []ReportAlert  `json:"alerts"`
	Note    string         `json:"note"`
}

type Report struct {
	BaseModel
	CowID       string        `json:"cow_id" gorm:"not null;index"`
	PeriodStart time.Time     `json:"period_start" gorm:"not null"`
	PeriodEnd   time.Time     `json:"period_end" gorm:"not null"`
	Summary     string        `json:"summary" gorm:"not null"`
	Score       float64       `json:"score" gorm:"not null"`
	Details     ReportDetails `json:"details_json" gorm:"column:details_json;type:jsonb;serializer:json;not null"`
}
