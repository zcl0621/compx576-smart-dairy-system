package model

import (
	"fmt"

	"gorm.io/gorm"
)

type CowStatus string

const (
	CowStatusInFarm   CowStatus = "in_farm"
	CowStatusSold     CowStatus = "sold"
	CowStatusInactive CowStatus = "inactive"
)

type CowCondition string

const (
	CowConditionNormal   CowCondition = "normal"
	CowConditionWarning  CowCondition = "warning"
	CowConditionCritical CowCondition = "critical"
	CowConditionOffline  CowCondition = "offline"
)

type Cow struct {
	BaseModel
	Name       string    `json:"name" gorm:"not null"`
	Tag        string    `json:"tag" gorm:"not null;uniqueIndex"`
	Age        int       `json:"age" gorm:"default:0"`
	CanMilking bool      `json:"can_milking" gorm:"default:false"`
	Status     CowStatus `json:"status" gorm:"default:in_farm"`
}

func OrderByCowCondition(db *gorm.DB) *gorm.DB {
	return db.Order(fmt.Sprintf(`
		CASE
			WHEN EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = '%s' AND alerts.severity = '%s' AND alerts.deleted_at IS NULL) THEN 1
			WHEN EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = '%s' AND alerts.severity = '%s' AND alerts.deleted_at IS NULL) THEN 2
			WHEN EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = '%s' AND alerts.severity = '%s' AND alerts.deleted_at IS NULL) THEN 4
			ELSE 3
		END`,
		AlertStatusActive, AlertSeverityCritical,
		AlertStatusActive, AlertSeverityWarning,
		AlertStatusActive, AlertSeverityOffline,
	))
}

// ConditionSubQuery builds the condition CASE expression for use in a SELECT list.
func ConditionSubQuery() string {
	return fmt.Sprintf(`(CASE
		WHEN EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = '%s' AND alerts.severity = '%s' AND alerts.deleted_at IS NULL) THEN '%s'
		WHEN EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = '%s' AND alerts.severity = '%s' AND alerts.deleted_at IS NULL) THEN '%s'
		WHEN EXISTS (SELECT 1 FROM alerts WHERE alerts.cow_id = cows.id AND alerts.status = '%s' AND alerts.severity = '%s' AND alerts.deleted_at IS NULL) THEN '%s'
		ELSE '%s'
	END)`,
		AlertStatusActive, AlertSeverityCritical, CowConditionCritical,
		AlertStatusActive, AlertSeverityWarning, CowConditionWarning,
		AlertStatusActive, AlertSeverityOffline, CowConditionOffline,
		CowConditionNormal,
	)
}
