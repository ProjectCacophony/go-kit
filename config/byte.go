package config

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func setByte(db *gorm.DB, namespace, key string, value []byte) error {
	var item Item

	return db.Where(&Item{
		Namespace: namespace,
		Key:       key,
	}).Assign(&Item{
		Value: value,
	}).FirstOrCreate(&item).Error
}

func getByte(db *gorm.DB, namespace, key string) ([]byte, error) {
	var item Item

	err := db.Where("namespace = ? AND key = ?", namespace, key).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []byte{}, nil
		}
		return nil, err
	}

	return item.Value, nil
}
