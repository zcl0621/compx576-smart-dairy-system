package model

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	if e := db.AutoMigrate(&User{}); e != nil {
		return e
	}
	if e := db.AutoMigrate(&Cow{}); e != nil {
		return e
	}
	if e := db.AutoMigrate(&Report{}); e != nil {
		return e
	}
	if e := db.AutoMigrate(&Metric{}); e != nil {
		return e
	}
	if e := db.AutoMigrate(&Alert{}); e != nil {
		return e
	}

	if db.Migrator().HasColumn(&Cow{}, "condition") {
		if e := db.Exec("ALTER TABLE cows DROP COLUMN condition").Error; e != nil {
			return e
		}
	}

	// partial unique index — AutoMigrate does not create partial indexes
	if e := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS alerts_cow_metric_active
		ON alerts (cow_id, metric_key)
		WHERE status = 'active' AND deleted_at IS NULL`).Error; e != nil {
		return e
	}

	return db.Exec(`CREATE INDEX IF NOT EXISTS metrics_cow_type_time
		ON metrics (cow_id, metric_type, created_at DESC)`).Error
}
