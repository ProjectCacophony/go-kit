package permissions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/interfaces"
)

const (
	// DiscordCacophonyServerGuildID is the Guild ID of the Cacophony Guild
	DiscordCacophonyServerGuildID = "435420687906111498"
)

type CacophonyBotPermission struct {
	name  string
	match func(state interfaces.State, userID string, channelID string, dm bool) bool
}

func newCacophonyBotAdmin() *CacophonyBotPermission {
	return &CacophonyBotPermission{
		name: "Bot Admin",
		match: func(
			state interfaces.State,
			userID string,
			channelID string,
			dm bool,
		) bool {
			aPermissions, err := state.UserPermissions(userID, DiscordCacophonyServerGuildID)
			if err != nil {
				return false
			}

			if aPermissions&discordgo.PermissionManageServer != discordgo.PermissionManageServer {
				return false
			}

			return true
		},
	}
}

func (p *CacophonyBotPermission) Name() string {
	return p.name
}

func (p *CacophonyBotPermission) Match(
	state interfaces.State,
	userID string,
	channelID string,
	dm bool,
) bool {
	return p.match(state, userID, channelID, dm)
}

var (
	// BotAdmin has Manage_Server permissions on the Bot Server
	BotAdmin = newCacophonyBotAdmin()
)
