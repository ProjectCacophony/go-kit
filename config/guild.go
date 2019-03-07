package config

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func guildNamespace(userID string) string {
	return "guild:" + userID
}

func GuildSetString(db *gorm.DB, guildID string, key, value string) error {
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
