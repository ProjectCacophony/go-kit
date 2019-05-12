package permissions

import "gitlab.com/Cacophony/go-kit/interfaces"

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
	userID string,
	channelID string,
	dm bool,
) bool {
	return dm
}

var (
	DiscordChannelDM = newDiscordDM()
)
