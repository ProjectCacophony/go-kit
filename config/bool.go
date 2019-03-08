package config

import (
	"github.com/jinzhu/gorm"
)

func setBool(db *gorm.DB, namespace, key string, value bool) error {
	var item Item

	var val byte
	if value {
		val = 1
	}

	return db.Where(&Item{
		Namespace: namespace,
		Key:       key,
	}).Assign(&Item{
		Value: []byte{val},
	}).FirstOrCreate(&item).Error
}

func getBool(db *gorm.DB, namespace, key string) (bool, error) {
	var item Item

	err := db.Where("namespace = ? AND key = ?", namespace, key).First(&item).Error
	if err != nil {
		return false, err
	}

	if len(item.Value) > 0 && item.Value[0] > 0 {
		return true, nil
	}

	return false, nil
}
