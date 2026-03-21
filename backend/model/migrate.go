package model

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	if e := db.AutoMigrate(&User{}); e != nil {
		return e
	}
	if e := db.AutoMigrate(Cow{}); e != nil {
		return e
	}
	if e := db.AutoMigrate(Report{}); e != nil {
		return e
	}
	if e := db.AutoMigrate(Metric{}); e != nil {
		return e
	}
	if e := db.AutoMigrate(Alert{}); e != nil {
		return e
	}
	return nil
}
