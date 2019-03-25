package config

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

func setInterface(db *gorm.DB, namespace, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return setByte(db, namespace, key, data)
}

func getInterface(db *gorm.DB, namespace, key string, value interface{}) error {
	data, err := getByte(db, namespace, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}
