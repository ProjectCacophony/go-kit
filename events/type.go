package events

// Type defines the type for a Event
type Type string

// defines various Event Types
const (
	ChannelCreateType            Type = "channel_create"
	ChannelDeleteType            Type = "channel_delete"
	ChannelPinsUpdateType        Type = "channel_pins_update"
	ChannelUpdateType            Type = "channel_update"
	GuildBanAddType              Type = "guild_ban_add"
	GuildBanRemoveType           Type = "guild_ban_remove"
	GuildCreateType              Type = "guild_create"
	GuildDeleteType              Type = "guild_delete"
	GuildEmojisUpdateType        Type = "guild_emojis_update"
	GuildMemberAddType           Type = "guild_member_add"
	GuildMemberRemoveType        Type = "guild_member_remove"
	GuildMemberUpdateType        Type = "guild_member_update"
	GuildMembersChunkType        Type = "guild_members_chunk"
	GuildRoleCreateType          Type = "guild_role_create"
	GuildRoleDeleteType          Type = "guild_role_delete"
	GuildRoleUpdateType          Type = "guild_role_update"
	GuildUpdateType              Type = "guild_update"
	GuildIntegrationsUpdateType  Type = "guild_integrations_update"
	MessageCreateType            Type = "message_create"
	MessageDeleteType            Type = "message_delete"
	MessageDeleteBulkType        Type = "message_delete_bulk"
	MessageReactionAddType       Type = "message_reaction_add"
	MessageReactionRemoveType    Type = "message_reaction_remove"
	MessageReactionRemoveAllType Type = "message_reaction_remove_all"
	MessageUpdateType            Type = "message_update"
	PresenceUpdateType           Type = "presence_update"
	UserUpdateType               Type = "user_update"
	VoiceStateUpdateType         Type = "voice_state_update"
	VoiceServerUpdateType        Type = "voice_server_update"
	WebhooksUpdateType           Type = "webhooks_update"
)
