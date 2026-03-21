package dashboard

import (
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

type SummaryResponse struct {
	TotalCows int64 `json:"total_cows"`
	Normal    int64 `json:"normal"`
	Warning   int64 `json:"warning"`
	Critical  int64 `json:"critical"`
	Offline   int64 `json:"offline"`
}

type ListItem struct {
	ID           string             `json:"id" gorm:"column:id"`
	Name         string             `json:"name" gorm:"column:name"`
	Tag          string             `json:"tag" gorm:"column:tag"`
	Condition    model.CowCondition `json:"condition" gorm:"column:condition"`
	Temperature  *float64           `json:"temperature" gorm:"column:temperature"`
	HeartRate    *float64           `json:"heart_rate" gorm:"column:heart_rate"`
	BloodOxygen  *float64           `json:"blood_oxygen" gorm:"column:blood_oxygen"`
	AlertMessage *string            `json:"alert_message" gorm:"column:alert_message"`
	UpdatedAt    time.Time          `json:"updated_at" gorm:"column:updated_at"`
}

type ListResponse struct {
	List []ListItem `json:"list"`
	common.PageResponse
}
