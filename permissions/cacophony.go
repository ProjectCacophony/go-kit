package permissions

import (
	"gitlab.com/Cacophony/go-kit/state"
)

type CacophonyBotOwner struct {
	name string
}

func newCacophonyBotOwner() *CacophonyBotOwner {
	return &CacophonyBotOwner{
		name: "Bot Owner",
	}
}

func (p *CacophonyBotOwner) Name() string {
	return p.name
}

func (p *CacophonyBotOwner) Match(
	state *state.State, botOwnerIDs []string, userID, channelID string, dm bool) bool {
	for _, botOwnerID := range botOwnerIDs {
		if botOwnerID != userID {
			continue
		}

		return true
	}
	return false
}

// nolint: gochecknoglobals
var (
	BotOwner = newCacophonyBotOwner()
)
