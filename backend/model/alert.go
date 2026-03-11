package model

type Alert struct {
	BaseModel
	CowID   string `json:"cow_id" gorm:"not null;index"`
	Message string `json:"message" gorm:"not null"`
}
