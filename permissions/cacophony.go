package permissions

import (
	"gitlab.com/Cacophony/go-kit/state"
)

// TODO (snake): change to bot permission? ex: bot owner, bot staff, bot mod
//     type CacophonyBotPermission
//     newBotPermission
//     ...

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
