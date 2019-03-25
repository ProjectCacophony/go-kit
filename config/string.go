package config

import (
	"github.com/jinzhu/gorm"
)

func setString(db *gorm.DB, namespace, key, value string) error {
	return setByte(db, namespace, key, []byte(value))
}

func getString(db *gorm.DB, namespace, key string) (string, error) {
	value, err := getByte(db, namespace, key)
	if err != nil {
		return "", err
	}

	return string(value), nil
}
