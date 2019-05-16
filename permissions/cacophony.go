package permissions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/interfaces"
)

const (
	// DiscordCacophonyServerGuildID is the Guild ID of the Cacophony Guild
	DiscordCacophonyServerGuildID = "435420687906111498"
	// DiscordCacophonyPatronRoleID is the Role ID of the Patron role on the Cacophony Guild
	DiscordCacophonyPatronRoleID = "578650002361155626"
)

type CacophonyBotPermission struct {
	name  string
	match func(state interfaces.State, userID string, channelID string, dm bool) bool
}

func newCacophonyBotAdmin(guildID string) *CacophonyBotPermission {
	return &CacophonyBotPermission{
		name: "Bot Admin",
		match: func(
			state interfaces.State,
			userID string,
			channelID string,
			dm bool,
		) bool {
			aPermissions, err := state.UserPermissions(userID, guildID)
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

func newCacophonyPatron(guildID, roleID string) *CacophonyBotPermission {
	return &CacophonyBotPermission{
		name: "Patron Supporter",
		match: func(state interfaces.State, userID string, channelID string, dm bool) bool {
			member, err := state.Member(guildID, userID)
			if err != nil {
				return false
			}

			for _, roleItemID := range member.Roles {
				if roleItemID == roleID {
					return true
				}
			}

			return false
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
	BotAdmin = newCacophonyBotAdmin(
		DiscordCacophonyServerGuildID,
	)
	// Patron has the Patron Role on the Bot Server
	Patron = newCacophonyPatron(
		DiscordCacophonyServerGuildID,
		DiscordCacophonyPatronRoleID,
	)
)
