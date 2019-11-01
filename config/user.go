package config

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func userNamespace(userID string) string {
	return "user:" + userID
}

func UserSetString(db *gorm.DB, userID, key, value string) error {
	if userID == "" {
		return errors.New("invalid User ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return setString(db, userNamespace(userID), key, value)
}

func UserGetString(db *gorm.DB, userID, key string) (string, error) {
	if userID == "" {
		return "", errors.New("invalid User ID")
	}
	if key == "" {
		return "", errors.New("invalid configuration key")
	}

	return getString(db, userNamespace(userID), key)
}

func UserSetBool(db *gorm.DB, userID, key string, value bool) error {
	if userID == "" {
		return errors.New("invalid User ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return setBool(db, userNamespace(userID), key, value)
}

func UserGetBool(db *gorm.DB, userID, key string) (bool, error) {
	if userID == "" {
		return false, errors.New("invalid User ID")
	}
	if key == "" {
		return false, errors.New("invalid configuration key")
	}

	return getBool(db, userNamespace(userID), key)
}

func UserSetInterface(db *gorm.DB, userID, key string, value interface{}) error {
	if userID == "" {
		return errors.New("invalid User ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return setInterface(db, userNamespace(userID), key, value)
}

func UserGetInterface(db *gorm.DB, userID, key string, value interface{}) error {
	if userID == "" {
		return errors.New("invalid User ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return getInterface(db, userNamespace(userID), key, value)
}

func UserSetBytes(db *gorm.DB, userID, key string, value []byte) error {
	if userID == "" {
		return errors.New("invalid User ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return setByte(db, userNamespace(userID), key, value)
}

func UserGetBytes(db *gorm.DB, userID, key string) ([]byte, error) {
	if userID == "" {
		return nil, errors.New("invalid User ID")
	}
	if key == "" {
		return nil, errors.New("invalid configuration key")
	}

	return getByte(db, userNamespace(userID), key)
}
