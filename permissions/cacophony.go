package permissions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/interfaces"
)

const (
	// DiscordCacophonyServerGuildID is the Guild ID of the Cacophony Guild
	DiscordCacophonyServerGuildID = "403003926720413699"
	// DiscordCacophonyPatronRoleID is the Role ID of the Patron role on the Cacophony Guild
	DiscordCacophonyPatronRoleID = "578650002361155626"
)

type CacophonyBotPermission struct {
	name  string
	match func(
		state interfaces.State,
		db *gorm.DB,
		userID string,
		channelID string,
		dm bool,
	) bool
}

func newCacophonyBotAdmin(guildID string) *CacophonyBotPermission {
	return &CacophonyBotPermission{
		name: "Bot Admin",
		match: func(
			state interfaces.State,
			db *gorm.DB,
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
		match: func(
			state interfaces.State,
			db *gorm.DB,
			userID string,
			channelID string,
			dm bool,
		) bool {
			// check Patron Role on Discord
			member, err := state.Member(guildID, userID)
			if err != nil {
				return false
			}

			for _, roleItemID := range member.Roles {
				if roleItemID == roleID {
					return true
				}
			}

			// check Patron DB
			var results int
			err = db.
				Table("patrons").
				Where(
					"discord_user_id = ? AND patron_status = ?",
					userID,
					"active_patron",
				).Count(&results).Error
			if err != nil {
				return false
			}

			return results > 0
		},
	}
}

func (p *CacophonyBotPermission) Name() string {
	return p.name
}

func (p *CacophonyBotPermission) Match(
	state interfaces.State,
	db *gorm.DB,
	userID string,
	channelID string,
	dm bool,
) bool {
	return p.match(
		state,
		db,
		userID,
		channelID,
		dm,
	)
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
