package model

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string         `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Index     int64          `json:"index" gorm:"autoIncrement"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	b.ID = xid.New().String()
	return nil
}
