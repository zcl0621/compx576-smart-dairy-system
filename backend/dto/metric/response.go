package metric

import (
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
	"github.com/zcl0621/compx576-smart-dairy-system/model"
)

type ListItem struct {
	ID          string             `json:"id" gorm:"column:id"`
	CowID       string             `json:"cow_id" gorm:"column:cow_id"`
	CowName     string             `json:"cow_name" gorm:"column:cow_name"`
	Source      model.MetricSource `json:"source" gorm:"column:source"`
	MetricType  model.MetricType   `json:"metric_type" gorm:"column:metric_type"`
	MetricValue float64            `json:"metric_value" gorm:"column:metric_value"`
	Unit        model.MetricUnit   `json:"unit" gorm:"column:unit"`
	CreatedAt   time.Time          `json:"created_at" gorm:"column:created_at"`
}

type ListResponse struct {
	List []ListItem `json:"list"`
	common.PageResponse
}
