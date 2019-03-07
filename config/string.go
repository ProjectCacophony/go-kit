package config

import (
	"github.com/jinzhu/gorm"
)

func setString(db *gorm.DB, namespace, key, value string) error {
	var item Item

	return db.Where(&Item{
		Namespace: namespace,
		Key:       key,
	}).Assign(&Item{
		Value: []byte(value),
	}).FirstOrCreate(&item).Error
}

func getString(db *gorm.DB, namespace, key string) (string, error) {
	var item Item

	err := db.Where("namespace = ? AND key = ?", namespace, key).First(&item).Error
	if err != nil {
		return "", err
	}

	return string(item.Value), nil
}
