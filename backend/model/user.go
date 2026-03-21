package model

type User struct {
	BaseModel
	Username string `json:"username" gorm:"not null;uniqueIndex"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"not null;uniqueIndex"`
}
