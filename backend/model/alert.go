package model

import "time"

type AlertSeverity string

const (
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
	AlertSeverityOffline  AlertSeverity = "offline"
)

type AlertStatus string

const (
	AlertStatusActive   AlertStatus = "active"
	AlertStatusResolved AlertStatus = "resolved"
)

type Alert struct {
	BaseModel
	CowID      string        `json:"cow_id" gorm:"not null;index"`
	ReportID   *string       `json:"report_id"`
	MetricKey  MetricType    `json:"metric_key" gorm:"not null"`
	Title      string        `json:"title" gorm:"not null"`
	Message    string        `json:"message" gorm:"not null"`
	Severity   AlertSeverity `json:"severity" gorm:"not null"`
	Status     AlertStatus   `json:"status" gorm:"not null;default:active"`
	ResolvedAt *time.Time    `json:"resolved_at"`
}
