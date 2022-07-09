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

	key += string(event.Status)

	if event.Since != nil {
		key += strconv.Itoa(*event.Since)
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
		strconv.Itoa(int(event.DefaultMessageNotifications)) +
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
		event.PremiumSince.String() +
		event.JoinedAt.String() +
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
		strconv.FormatInt(event.Permissions, 10) +
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
		strconv.Itoa(int(event.Type)) +
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
	return event.ID + strconv.Itoa(int(event.Type)) +
		strconv.FormatInt(event.Allow, 10) + strconv.FormatInt(event.Deny, 10)
}

func messageUpdateKey(event *discordgo.MessageUpdate) string {
	return event.ID + event.EditedTimestamp.String()
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
