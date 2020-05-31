package permissions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/interfaces"
)

type Discord struct {
	name string
	id   int
}

func NewDiscordPermission(name string, id int) *Discord {
	return &Discord{
		name: name,
		id:   id,
	}
}

func (p *Discord) Name() string {
	return p.name
}

func (p *Discord) Match(
	state interfaces.State,
	db *gorm.DB,
	userID string,
	channelID string,
	dm bool,
	superUser bool,
) bool {
	if dm {
		return false
	}

	if superUser {
		return true
	}

	return discord.UserHasPermission(state, userID, channelID, p.id)
}

var (
	DiscordReadMessages = NewDiscordPermission(
		"Read Messages",
		discordgo.PermissionReadMessages,
	)
	DiscordSendMessages = NewDiscordPermission(
		"Send Messages",
		discordgo.PermissionSendMessages,
	)
	DiscordSendTTSMessages = NewDiscordPermission(
		"Send TTS Messages",
		discordgo.PermissionSendTTSMessages,
	)
	DiscordManageMessages = NewDiscordPermission(
		"Manage Messages",
		discordgo.PermissionManageMessages,
	)
	DiscordEmbedLinks = NewDiscordPermission(
		"Embed Links",
		discordgo.PermissionEmbedLinks,
	)
	DiscordAttachFiles = NewDiscordPermission(
		"Attach Files",
		discordgo.PermissionAttachFiles,
	)
	DiscordReadMessageHistory = NewDiscordPermission(
		"Read Message History",
		discordgo.PermissionReadMessageHistory,
	)
	DiscordMentionEveryone = NewDiscordPermission(
		"Mention Everyone",
		discordgo.PermissionMentionEveryone,
	)
	DiscordUseExternalEmojis = NewDiscordPermission(
		"Use External Emojis",
		discordgo.PermissionUseExternalEmojis,
	)

	DiscordVoiceConnect = NewDiscordPermission(
		"Voice Connect",
		discordgo.PermissionVoiceConnect,
	)
	DiscordVoiceSpeak = NewDiscordPermission(
		"Voice Speak",
		discordgo.PermissionVoiceSpeak,
	)
	DiscordVoiceMuteMembers = NewDiscordPermission(
		"Voice Mute Members",
		discordgo.PermissionVoiceMuteMembers,
	)
	DiscordVoiceDeafenMembers = NewDiscordPermission(
		"Voice Deafen Members",
		discordgo.PermissionVoiceDeafenMembers,
	)
	DiscordVoiceMoveMembers = NewDiscordPermission(
		"Voice Move Members",
		discordgo.PermissionVoiceMoveMembers,
	)
	DiscordVoiceUseVAD = NewDiscordPermission(
		"Voice Use VAD",
		discordgo.PermissionVoiceUseVAD,
	)

	DiscordChangeNickname = NewDiscordPermission(
		"Change Nickname",
		discordgo.PermissionChangeNickname,
	)
	DiscordManageNicknames = NewDiscordPermission(
		"Manage Nicknames",
		discordgo.PermissionManageNicknames,
	)
	DiscordManageRoles = NewDiscordPermission(
		"Manage Roles",
		discordgo.PermissionManageRoles,
	)
	DiscordManageWebhooks = NewDiscordPermission(
		"Manage Webhooks",
		discordgo.PermissionManageWebhooks,
	)
	DiscordManageEmojis = NewDiscordPermission(
		"Manage Emojis",
		discordgo.PermissionManageEmojis,
	)

	DiscordCreateInstantInvite = NewDiscordPermission(
		"Create Instant Invite",
		discordgo.PermissionCreateInstantInvite,
	)
	DiscordKickMembers = NewDiscordPermission(
		"Kick Members",
		discordgo.PermissionKickMembers,
	)
	DiscordBanMembers = NewDiscordPermission(
		"Ban Members",
		discordgo.PermissionBanMembers,
	)
	DiscordAdministrator = NewDiscordPermission(
		"Administrator",
		discordgo.PermissionAdministrator,
	)
	DiscordManageChannels = NewDiscordPermission(
		"Manage Channels",
		discordgo.PermissionManageChannels,
	)
	DiscordManageServer = NewDiscordPermission(
		"Manage Server",
		discordgo.PermissionManageServer,
	)
	DiscordAddReactions = NewDiscordPermission(
		"Add Reactions",
		discordgo.PermissionAddReactions,
	)
	DiscordViewAuditLogs = NewDiscordPermission(
		"View Audit Logs",
		discordgo.PermissionViewAuditLogs,
	)
)
