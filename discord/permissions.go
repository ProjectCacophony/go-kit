package discord

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/state"
)

// UserHasPermission returns true if the User has all of th egiven permissions in the given channel
func UserHasPermission(
	state *state.State, userID, channelID string, firstPermission int, permissions ...int,
) bool {
	if userID == "" || channelID == "" {
		return false
	}

	if state == nil {
		return false
	}

	userChannelPermissions, err := state.UserChannelPermissions(userID, channelID)
	if err != nil {
		return false
	}

	for _, permission := range append(permissions, firstPermission) {

		if userChannelPermissions&permission != permission {
			return false
		}
	}

	return true
}

// UserHasPermissionOr returns true if the User has any of the given permissions in the given channel
func UserHasPermissionOr(
	state *state.State, userID, channelID string, firstPermission int, permissions ...int,
) bool {
	if userID == "" || channelID == "" {
		return false
	}

	if state == nil {
		return false
	}

	userChannelPermissions, err := state.UserChannelPermissions(userID, channelID)
	if err != nil {
		return false
	}

	for _, permission := range append(permissions, firstPermission) {

		if userChannelPermissions&permission == permission {
			return true
		}
	}

	return false
}

// PermissionName returns a human readable name for a permission
func PermissionName(permission int) string {
	var name string

	switch permission {
	case discordgo.PermissionReadMessages:
		name = "Read Messages"
	case discordgo.PermissionSendMessages:
		name = "Send Messages"
	case discordgo.PermissionSendTTSMessages:
		name = "Send TTS Messages"
	case discordgo.PermissionManageMessages:
		name = "Manage Messages"
	case discordgo.PermissionEmbedLinks:
		name = "Embed Links"
	case discordgo.PermissionAttachFiles:
		name = "Attach Files"
	case discordgo.PermissionReadMessageHistory:
		name = "Read Message History"
	case discordgo.PermissionMentionEveryone:
		name = "Mention Everyone"
	case discordgo.PermissionUseExternalEmojis:
		name = "Use External Emojis"

	case discordgo.PermissionVoiceConnect:
		name = "Voice Connect"
	case discordgo.PermissionVoiceSpeak:
		name = "Voice Speak"
	case discordgo.PermissionVoiceMuteMembers:
		name = "Mute Members"
	case discordgo.PermissionVoiceDeafenMembers:
		name = "Voice Deafen Members"
	case discordgo.PermissionVoiceMoveMembers:
		name = "Voice Move Members"
	case discordgo.PermissionVoiceUseVAD:
		name = "Voice Use VAD"

	case discordgo.PermissionChangeNickname:
		name = "Change Nicknames"
	case discordgo.PermissionManageNicknames:
		name = "Manage Nicknames"
	case discordgo.PermissionManageRoles:
		name = "Manage Roles"
	case discordgo.PermissionManageWebhooks:
		name = "Manage Webhooks"
	case discordgo.PermissionManageEmojis:
		name = "Manage Emojis"

	case discordgo.PermissionCreateInstantInvite:
		name = "Create Instant Invite"
	case discordgo.PermissionKickMembers:
		name = "Kick Members"
	case discordgo.PermissionBanMembers:
		name = "Ban Members"
	case discordgo.PermissionAdministrator:
		name = "Administrator"
	case discordgo.PermissionManageChannels:
		name = "Manage Channels"
	case discordgo.PermissionManageServer:
		name = "Manage Server"
	case discordgo.PermissionAddReactions:
		name = "Add Reactions"
	case discordgo.PermissionViewAuditLogs:
		name = "View Audit Logs"
	}

	return name
}
