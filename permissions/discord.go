package permissions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/state"
)

type Discord struct {
	name string
	id   int
}

func newDiscordPermission(name string, id int) *Discord {
	return &Discord{
		name: name,
		id:   id,
	}
}

func (p *Discord) Name() string {
	return p.name
}

func (p *Discord) Match(state *state.State, botOwnerIDs []string, userID, channelID string) bool {
	return discord.UserHasPermission(state, userID, channelID, p.id)
}

// nolint: gochecknoglobals
var (
	DiscordAdministrator = newDiscordPermission(
		"Administrator",
		discordgo.PermissionAdministrator,
	)
	DiscordManageChannels = newDiscordPermission(
		"Manage Channels",
		discordgo.PermissionManageChannels,
	)
)

// TODO: add all permissions
