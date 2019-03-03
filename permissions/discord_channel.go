package permissions

import (
	"gitlab.com/Cacophony/go-kit/state"
)

type DiscordDM struct {
	yes bool
}

func newDiscordDM(yes bool) *DiscordDM {
	return &DiscordDM{
		yes: yes,
	}
}

func (p *DiscordDM) Name() string {
	if p.yes {
		return "DM"
	}

	return "not in DM"
}

func (p *DiscordDM) Match(state *state.State, botOwnerIDs []string, userID, channelID string, dm bool) bool {
	return p.yes == dm
}

// nolint: gochecknoglobals
var (
	DiscordChannelDM    = newDiscordDM(true)
	DiscordChannelNotDM = newDiscordDM(false)
)
