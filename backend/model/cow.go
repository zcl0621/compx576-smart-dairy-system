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
	Name       string       `json:"name" gorm:"not null"`
	Tag        string       `json:"tag" gorm:"not null;uniqueIndex"`
	Age        int          `json:"age" gorm:"default:0"`
	CanMilking bool         `json:"can_milking" gorm:"default:false"`
	Status     CowStatus    `json:"status" gorm:"default:in_farm"`
	Condition  CowCondition `json:"condition" gorm:"default:normal"`
}

func OrderByCowCondition(db *gorm.DB) *gorm.DB {
	return db.Order(fmt.Sprintf(`
		CASE condition
			WHEN '%s' THEN 1
			WHEN '%s' THEN 2
			WHEN '%s' THEN 3
			WHEN '%s' THEN 4
			ELSE 999
		END
	`, CowConditionCritical, CowConditionWarning, CowConditionNormal, CowConditionOffline))
}
