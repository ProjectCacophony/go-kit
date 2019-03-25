package config

import (
	"github.com/jinzhu/gorm"
)

func setBool(db *gorm.DB, namespace, key string, value bool) error {
	var val byte
	if value {
		val = 1
	}

	return setByte(db, namespace, key, []byte{val})
}

func getBool(db *gorm.DB, namespace, key string) (bool, error) {
	value, err := getByte(db, namespace, key)
	if err != nil {
		return false, err
	}

	if len(value) > 0 && value[0] > 0 {
		return true, nil
	}

	return false, nil
}
