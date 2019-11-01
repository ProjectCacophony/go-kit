package config

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func guildNamespace(userID string) string {
	return "guild:" + userID
}

func GuildSetString(db *gorm.DB, guildID, key, value string) error {
	if guildID == "" {
		return errors.New("invalid Guild ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return setString(db, guildNamespace(guildID), key, value)
}

func GuildGetString(db *gorm.DB, guildID, key string) (string, error) {
	if guildID == "" {
		return "", errors.New("invalid Guild ID")
	}
	if key == "" {
		return "", errors.New("invalid configuration key")
	}

	return getString(db, guildNamespace(guildID), key)
}

func GuildSetBool(db *gorm.DB, guildID, key string, value bool) error {
	if guildID == "" {
		return errors.New("invalid Guild ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return setBool(db, guildNamespace(guildID), key, value)
}

func GuildGetBool(db *gorm.DB, guildID, key string) (bool, error) {
	if guildID == "" {
		return false, errors.New("invalid Guild ID")
	}
	if key == "" {
		return false, errors.New("invalid configuration key")
	}

	return getBool(db, guildNamespace(guildID), key)
}

func GuildSetInterface(db *gorm.DB, guildID, key string, value interface{}) error {
	if guildID == "" {
		return errors.New("invalid Guild ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return setInterface(db, guildNamespace(guildID), key, value)
}

func GuildGetInterface(db *gorm.DB, guildID, key string, value interface{}) error {
	if guildID == "" {
		return errors.New("invalid Guild ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return getInterface(db, guildNamespace(guildID), key, value)
}

func GuildSetBytes(db *gorm.DB, guildID, key string, value []byte) error {
	if guildID == "" {
		return errors.New("invalid Guild ID")
	}
	if key == "" {
		return errors.New("invalid configuration key")
	}

	return setByte(db, guildNamespace(guildID), key, value)
}

func GuildGetBytes(db *gorm.DB, guildID, key string) ([]byte, error) {
	if guildID == "" {
		return nil, errors.New("invalid Guild ID")
	}
	if key == "" {
		return nil, errors.New("invalid configuration key")
	}

	return getByte(db, guildNamespace(guildID), key)
}
