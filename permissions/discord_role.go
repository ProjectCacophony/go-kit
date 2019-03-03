package permissions

import (
	"gitlab.com/Cacophony/go-kit/state"
)

type DiscordRole struct {
	guildID string
	roleID  string
	state   *state.State
}

func NewDiscordRole(state *state.State, guildID, roleID string) *DiscordRole {
	return &DiscordRole{
		guildID: guildID,
		roleID:  roleID,
		state:   state,
	}
}

func (p *DiscordRole) Name() string {
	role, err := p.state.Role(p.guildID, p.roleID)
	if err != nil {
		return "@#" + p.roleID
	}

	var guildPrefix string
	guild, err := p.state.Guild(p.guildID)
	if err == nil {
		guildPrefix = guild.Name + " / "
	}

	return guildPrefix + "@" + role.Name
}

func (p *DiscordRole) Match(state *state.State, botOwnerIDs []string, userID, channelID string, dm bool) bool {
	if dm {
		return false
	}

	member, err := state.Member(p.guildID, userID)
	if err != nil {
		return false
	}

	for _, roleID := range member.Roles {
		if roleID != p.roleID {
			continue
		}

		return true
	}

	return false
}
