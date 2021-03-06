package permissions

import (
	"errors"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/config"
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
	key   string
	match func(
		state interfaces.State,
		db *gorm.DB,
		userID string,
		channelID string,
		dm bool,
	) bool
	give func(
		db *gorm.DB,
		userID string,
		permission *CacophonyBotPermission,
	) error
	remove func(
		db *gorm.DB,
		userID string,
		permission *CacophonyBotPermission,
	) error
}

func newCacophonyBotAdmin(guildID string) *CacophonyBotPermission {
	userIDs := strings.Split(os.Getenv("BOT_ADMIN_USER_IDS"), ",")

	return &CacophonyBotPermission{
		name: "Bot Admin",
		match: func(
			state interfaces.State,
			db *gorm.DB,
			userID string,
			channelID string,
			dm bool,
		) bool {
			for _, adminUserID := range userIDs {
				if userID == adminUserID {
					return true
				}
			}

			aPermissions, err := state.UserPermissions(userID, guildID)
			if err != nil {
				return false
			}

			if aPermissions&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator {
				return false
			}

			return true
		},
		give: func(db *gorm.DB, userID string, permission *CacophonyBotPermission) error {
			return errors.New("bot admin permission cannot be set")
		},
		remove: func(db *gorm.DB, userID string, permission *CacophonyBotPermission) error {
			return errors.New("bot admin permission cannot be set")
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
			// bot admins have patron privileges
			if BotAdmin != nil && BotAdmin.Match(state, db, userID, channelID, dm, false) {
				return true
			}

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
		give: func(db *gorm.DB, userID string, permission *CacophonyBotPermission) error {
			// TODO
			return errors.New("todo")
		},
		remove: func(db *gorm.DB, userID string, permission *CacophonyBotPermission) error {
			// TODO
			return errors.New("todo")
		},
	}
}

/* Just commenting this out cause of the linter
func newCacophonyPermission(name, key string) *CacophonyBotPermission {
	return &CacophonyBotPermission{
		name: name,
		key:  key,
		match: func(
			state interfaces.State,
			db *gorm.DB,
			userID string,
			channelID string,
			dm bool,
		) bool {
			hasPermission, err := config.UserGetBool(db, userID, key)
			if err != nil {
				return false
			}

			return hasPermission
		},
		give: func(db *gorm.DB, userID string, permission *CacophonyBotPermission) error {
			return config.UserSetBool(db, userID, permission.key, true)
		},
		remove: func(db *gorm.DB, userID string, permission *CacophonyBotPermission) error {
			return config.UserSetBool(db, userID, permission.key, false)
		},
	}
}*/

// newCacophonyDefaultPermission is used when the given permission is defaulted to true
//   for example: file upload permissions should default to true unless revoked
func newCacophonyDefaultPermission(name, key string) *CacophonyBotPermission {
	return &CacophonyBotPermission{
		name: name,
		key:  key,
		match: func(
			state interfaces.State,
			db *gorm.DB,
			userID string,
			channelID string,
			dm bool,
		) bool {
			hasPermission, err := config.UserGetBool(db, userID, key)
			if err != nil {
				return true
			}

			return !hasPermission
		},
		give: func(db *gorm.DB, userID string, permission *CacophonyBotPermission) error {
			return config.UserSetBool(db, userID, permission.key, false)
		},
		remove: func(db *gorm.DB, userID string, permission *CacophonyBotPermission) error {
			return config.UserSetBool(db, userID, permission.key, true)
		},
	}
}

func (p *CacophonyBotPermission) Name() string {
	return p.name
}

func (p *CacophonyBotPermission) Match(state interfaces.State, db *gorm.DB, userID string, channelID string, dm bool, superUser bool) bool {
	return p.match(
		state,
		db,
		userID,
		channelID,
		dm,
	)
}

func (p *CacophonyBotPermission) Give(db *gorm.DB, userID string) error {
	return p.give(db, userID, p)
}

func (p *CacophonyBotPermission) Remove(db *gorm.DB, userID string) error {
	return p.remove(db, userID, p)
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

	// CacoFileStorage has ability to use commands that store files
	CacoFileStorage = newCacophonyDefaultPermission(
		"CacoFileStorage",
		"cacophony:permission:storage",
	)
)
