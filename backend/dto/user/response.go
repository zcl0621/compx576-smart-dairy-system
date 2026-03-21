package user

import (
	"time"

	"github.com/zcl0621/compx576-smart-dairy-system/dto/common"
)

type ListItem struct {
	ID        string    `json:"id" gorm:"column:id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	Username  string    `json:"username" gorm:"column:username"`
	Email     string    `json:"email" gorm:"column:email"`
}

type ListResponse struct {
	List []ListItem `json:"list"`
	common.PageResponse
}

type InfoResponse struct {
	ID        string    `json:"id" gorm:"column:id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	Username  string    `json:"username" gorm:"column:username"`
	Email     string    `json:"email" gorm:"column:email"`
}
