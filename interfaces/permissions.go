package interfaces

import (
	"github.com/jinzhu/gorm"
)

type Permission interface {
	Name() string
	Match(
		state State,
		db *gorm.DB,
		userID string,
		channelID string,
		dm bool,
		superUser bool,
	) bool
}
