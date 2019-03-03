package permissions

import (
	"gitlab.com/Cacophony/go-kit/state"
)

type DiscordDM struct {
}

func newDiscordDM() *DiscordDM {
	return &DiscordDM{}
}

func (p *DiscordDM) Name() string {
	return "DM"
}

func (p *DiscordDM) Match(state *state.State, botOwnerIDs []string, userID, channelID string, dm bool) bool {
	return dm
}

// nolint: gochecknoglobals
var (
	DiscordChannelDM = newDiscordDM()
)
