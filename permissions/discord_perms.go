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

func (p *Discord) Match(state *state.State, botOwnerIDs []string, userID, channelID string, dm bool) bool {
	if dm {
		return false
	}

	return discord.UserHasPermission(state, userID, channelID, p.id)
}

// nolint: gochecknoglobals
var (
	DiscordReadMessages = newDiscordPermission(
		"Read Messages",
		discordgo.PermissionReadMessages,
	)
	DiscordSendMessages = newDiscordPermission(
		"Send Messages",
		discordgo.PermissionSendMessages,
	)
	DiscordSendTTSMessages = newDiscordPermission(
		"Send TTS Messages",
		discordgo.PermissionSendTTSMessages,
	)
	DiscordManageMessages = newDiscordPermission(
		"Manage Messages",
		discordgo.PermissionManageMessages,
	)
	DiscordEmbedLinks = newDiscordPermission(
		"Embed Links",
		discordgo.PermissionEmbedLinks,
	)
	DiscordAttachFiles = newDiscordPermission(
		"Attach Files",
		discordgo.PermissionAttachFiles,
	)
	DiscordReadMessageHistory = newDiscordPermission(
		"Read Message History",
		discordgo.PermissionReadMessageHistory,
	)
	DiscordMentionEveryone = newDiscordPermission(
		"Mention Everyone",
		discordgo.PermissionMentionEveryone,
	)
	DiscordUseExternalEmojis = newDiscordPermission(
		"Use External Emojis",
		discordgo.PermissionUseExternalEmojis,
	)

	DiscordVoiceConnect = newDiscordPermission(
		"Voice Connect",
		discordgo.PermissionVoiceConnect,
	)
	DiscordVoiceSpeak = newDiscordPermission(
		"Voice Speak",
		discordgo.PermissionVoiceSpeak,
	)
	DiscordVoiceMuteMembers = newDiscordPermission(
		"Voice Mute Members",
		discordgo.PermissionVoiceMuteMembers,
	)
	DiscordVoiceDeafenMembers = newDiscordPermission(
		"Voice Deafen Members",
		discordgo.PermissionVoiceDeafenMembers,
	)
	DiscordVoiceMoveMembers = newDiscordPermission(
		"Voice Move Members",
		discordgo.PermissionVoiceMoveMembers,
	)
	DiscordVoiceUseVAD = newDiscordPermission(
		"Voice Use VAD",
		discordgo.PermissionVoiceUseVAD,
	)

	DiscordChangeNickname = newDiscordPermission(
		"Change Nickname",
		discordgo.PermissionChangeNickname,
	)
	DiscordManageNicknames = newDiscordPermission(
		"Manage Nicknames",
		discordgo.PermissionManageNicknames,
	)
	DiscordManageRoles = newDiscordPermission(
		"Manage Roles",
		discordgo.PermissionManageRoles,
	)
	DiscordManageWebhooks = newDiscordPermission(
		"Manage Webhooks",
		discordgo.PermissionManageWebhooks,
	)
	DiscordManageEmojis = newDiscordPermission(
		"Manage Emojis",
		discordgo.PermissionManageEmojis,
	)

	DiscordCreateInstantInvite = newDiscordPermission(
		"Create Instant Invite",
		discordgo.PermissionCreateInstantInvite,
	)
	DiscordKickMembers = newDiscordPermission(
		"Kick Members",
		discordgo.PermissionKickMembers,
	)
	DiscordBanMembers = newDiscordPermission(
		"Ban Members",
		discordgo.PermissionBanMembers,
	)
	DiscordAdministrator = newDiscordPermission(
		"Administrator",
		discordgo.PermissionAdministrator,
	)
	DiscordManageChannels = newDiscordPermission(
		"Manage Channels",
		discordgo.PermissionManageChannels,
	)
	DiscordManageServer = newDiscordPermission(
		"Manage Server",
		discordgo.PermissionManageServer,
	)
	DiscordAddReactions = newDiscordPermission(
		"Add Reactions",
		discordgo.PermissionAddReactions,
	)
	DiscordViewAuditLogs = newDiscordPermission(
		"View Audit Logs",
		discordgo.PermissionViewAuditLogs,
	)
)
