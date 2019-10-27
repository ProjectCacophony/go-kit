package events

// Type defines the type for a Event
type Type string

// defines various Event Types
const (
	// Discord Events
	ChannelCreateType            Type = "discord_channel_create"
	ChannelDeleteType            Type = "discord_channel_delete"
	ChannelPinsUpdateType        Type = "discord_channel_pins_update"
	ChannelUpdateType            Type = "discord_channel_update"
	GuildBanAddType              Type = "discord_guild_ban_add"
	GuildBanRemoveType           Type = "discord_guild_ban_remove"
	GuildCreateType              Type = "discord_guild_create"
	GuildDeleteType              Type = "discord_guild_delete"
	GuildEmojisUpdateType        Type = "discord_guild_emojis_update"
	GuildMemberAddType           Type = "discord_guild_member_add"
	GuildMemberRemoveType        Type = "discord_guild_member_remove"
	GuildMemberUpdateType        Type = "discord_guild_member_update"
	GuildMembersChunkType        Type = "discord_guild_members_chunk"
	GuildRoleCreateType          Type = "discord_guild_role_create"
	GuildRoleDeleteType          Type = "discord_guild_role_delete"
	GuildRoleUpdateType          Type = "discord_guild_role_update"
	GuildUpdateType              Type = "discord_guild_update"
	GuildIntegrationsUpdateType  Type = "discord_guild_integrations_update"
	MessageCreateType            Type = "discord_message_create"
	MessageDeleteType            Type = "discord_message_delete"
	MessageDeleteBulkType        Type = "discord_message_delete_bulk"
	MessageReactionAddType       Type = "discord_message_reaction_add"
	MessageReactionRemoveType    Type = "discord_message_reaction_remove"
	MessageReactionRemoveAllType Type = "discord_message_reaction_remove_all"
	MessageUpdateType            Type = "discord_message_update"
	PresenceUpdateType           Type = "discord_presence_update"
	UserUpdateType               Type = "discord_user_update"
	VoiceStateUpdateType         Type = "discord_voice_state_update"
	VoiceServerUpdateType        Type = "discord_voice_server_update"
	WebhooksUpdateType           Type = "discord_webhooks_update"
	// Cacophony Events
	CacophonyBucketUpdate           Type = "cacophony_bucket_update"
	CacophonyServerlistServerExpire Type = "cacophony_serverlist_server_expire"
	CacophonyQuickactionRemind      Type = "cacophony_quickaction_remind"
	CacophonyQuestionnaireMatch     Type = "cacophony_questionnaire_match"
	CacophonyAutomodWait            Type = "cacophony_automod_wait"
	CacophonyEventlogUpdate         Type = "cacophony_eventlog_update"
	CacophonyDiffGuild              Type = "cacophony_diff_guild"
	CacophonyDiffMember             Type = "cacophony_diff_member"
	CacophonyDiffChannel            Type = "cacophony_diff_channel"
	CacophonyDiffRole               Type = "cacophony_diff_role"
	CacophonyDiffEmoji              Type = "cacophony_diff_emoji"
	CacophonyDiffWebhooks           Type = "cacophony_diff_webhooks"
	CacophonyDiffInvites            Type = "cacophony_diff_invites"
	CacophonyGuildMemberAddExtra    Type = "cacophony_extra_guild_member_add"
)
