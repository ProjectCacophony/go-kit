package interfaces

import (
	"github.com/bwmarrin/discordgo"
)

type State interface {
	Channel(channelID string) (channel *discordgo.Channel, err error)
	Guild(guildID string) (guild *discordgo.Guild, err error)
	Member(guildID, userID string) (member *discordgo.Member, err error)
	Role(guildID, roleID string) (role *discordgo.Role, err error)

	UserPermissions(userID, guildID string) (apermissions int, err error)
	UserChannelPermissions(userID, channelID string) (apermissions int, err error)
}
