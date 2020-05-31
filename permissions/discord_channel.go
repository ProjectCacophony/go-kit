package permissions

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/interfaces"
)

type DiscordDM struct {
}

func newDiscordDM() *DiscordDM {
	return &DiscordDM{}
}

func (p *DiscordDM) Name() string {
	return "DM"
}

func (p *DiscordDM) Match(
	state interfaces.State,
	db *gorm.DB,
	userID string,
	channelID string,
	dm bool,
	superUser bool,
) bool {
	return dm
}

var (
	DiscordChannelDM = newDiscordDM()
)
