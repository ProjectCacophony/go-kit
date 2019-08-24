package events

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func guildMemberChunkKey(event *discordgo.GuildMembersChunk) string {
	key := event.GuildID
	for _, member := range event.Members {
		if member == nil || member.User == nil {
			continue
		}

		key += member.User.ID
	}

	return key
}

func presenceUpdateKey(event *discordgo.PresenceUpdate) string {
	key := event.GuildID

	if event.User != nil {
		key += event.User.ID
		key += event.User.Username
		key += event.User.Discriminator
		key += event.User.Avatar
	}

	key += event.Nick
	key += string(event.Status)
	key += strings.Join(event.Roles, "")

	if event.Game != nil {
		key += event.Game.State
		key += strconv.Itoa(int(event.Game.Type))
		key += event.Game.Name
		key += event.Game.URL
		key += event.Game.ApplicationID
		key += event.Game.Details
	}

	return key
}

func guildKey(event *discordgo.Guild) string {
	return event.ID +
		event.Name +
		event.Icon +
		event.Description +
		event.VanityURLCode +
		event.OwnerID +
		event.Banner +
		event.Splash +
		event.Region +
		event.AfkChannelID +
		strconv.Itoa(event.AfkTimeout) +
		strconv.Itoa(event.DefaultMessageNotifications) +
		event.EmbedChannelID +
		strconv.FormatBool(event.EmbedEnabled) +
		strconv.Itoa(int(event.ExplicitContentFilter)) +
		strings.Join(event.Features, "") +
		strconv.Itoa(int(event.MfaLevel)) +
		strconv.Itoa(int(event.PremiumTier)) +
		event.SystemChannelID +
		strconv.Itoa(int(event.VerificationLevel)) +
		event.WidgetChannelID +
		strconv.FormatBool(event.WidgetEnabled)
}

func memberKey(event *discordgo.Member) string {
	key := event.GuildID +
		strings.Join(event.Roles, "") +
		event.Nick +
		string(event.PremiumSince) +
		string(event.JoinedAt) +
		strconv.FormatBool(event.Deaf) +
		strconv.FormatBool(event.Mute)

	if event.User != nil {
		key += userKey(event.User)
	}

	return key
}

func userKey(event *discordgo.User) string {
	return event.ID +
		event.Avatar +
		event.Discriminator +
		event.Username +
		strconv.FormatBool(event.Bot) +
		event.Token +
		event.Email +
		event.Locale +
		strconv.FormatBool(event.MFAEnabled) +
		strconv.FormatBool(event.Verified)
}

func roleKey(event *discordgo.Role) string {
	return event.ID +
		strconv.Itoa(event.Color) +
		event.Name +
		strconv.Itoa(event.Permissions) +
		strconv.FormatBool(event.Hoist) +
		strconv.FormatBool(event.Managed) +
		strconv.FormatBool(event.Mentionable) +
		strconv.Itoa(event.Position)

}

func emojisUpdateKey(event *discordgo.GuildEmojisUpdate) string {
	key := event.GuildID

	for _, emoji := range event.Emojis {
		if emoji == nil {
			continue
		}

		key += emoji.ID +
			emoji.Name +
			strings.Join(emoji.Roles, "") +
			strconv.FormatBool(emoji.Managed) +
			strconv.FormatBool(emoji.Animated) +
			strconv.FormatBool(emoji.RequireColons)
	}

	return key
}

func channelKey(event *discordgo.Channel) string {
	key := event.ID +
		event.Name +
		event.GuildID +
		strconv.Itoa(event.Position) +
		string(event.Type) +
		event.ParentID +
		event.Icon +
		strconv.Itoa(event.Bitrate) +
		strconv.FormatBool(event.NSFW) +
		event.Topic +
		strconv.Itoa(event.UserLimit)

	for _, permissionOverwrite := range event.PermissionOverwrites {
		if permissionOverwrite == nil {
			continue
		}

		key += permissionOverwriteKey(permissionOverwrite)
	}

	for _, recipient := range event.Recipients {
		if recipient == nil {
			continue
		}

		key += userKey(recipient)
	}

	return key
}

func permissionOverwriteKey(event *discordgo.PermissionOverwrite) string {
	return event.ID +
		event.Type +
		strconv.Itoa(event.Allow) +
		strconv.Itoa(event.Deny)
}

func messageUpdateKey(event *discordgo.MessageUpdate) string {
	return event.ID +
		string(event.EditedTimestamp)
}

func messageReactionKey(event *discordgo.MessageReaction) string {
	return event.GuildID +
		event.ChannelID +
		event.UserID +
		event.MessageID +
		event.Emoji.ID +
		event.Emoji.Name
}

func voiceStateKey(event *discordgo.VoiceState) string {
	return event.GuildID +
		event.ChannelID +
		event.UserID +
		event.SessionID +
		strconv.FormatBool(event.Deaf) +
		strconv.FormatBool(event.Mute) +
		strconv.FormatBool(event.SelfDeaf) +
		strconv.FormatBool(event.SelfMute) +
		strconv.FormatBool(event.Suppress)
}