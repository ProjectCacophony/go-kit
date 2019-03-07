package config

import (
	"github.com/jinzhu/gorm"
)

func InitConfig(db *gorm.DB) error {
	return db.AutoMigrate(Item{}).Error
}
