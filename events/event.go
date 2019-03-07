package events

import (
	"context"
	"time"

	"gitlab.com/Cacophony/go-kit/discord"

	"github.com/go-redis/redis"

	"gitlab.com/Cacophony/go-kit/paginator"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

// Event represents an Event
// TODO: event ID or Key
type Event struct {
	ID         string    `json:"id"`
	CacheKey   string    `json:"cache_key"`
	Type       Type      `json:"type"`
	ReceivedAt time.Time `json:"received_at" dynamodbav:"received_at,unixtime"`
	BotUserID  string    `json:"bot_user_id"`
	GuildID    string    `json:"guild_id,omitempty"`
	ChannelID  string    `json:"channel_id,omitempty"`
	UserID     string    `json:"user_id,omitempty"`

	// discordgo event data
	ChannelCreate            *discordgo.ChannelCreate            `json:"discord_channel_create,omitempty" dynamodbav:"discord_channel_create,binaryset"`
	ChannelDelete            *discordgo.ChannelDelete            `json:"discord_channel_delete,omitempty" dynamodbav:"discord_channel_delete,binaryset"`
	ChannelPinsUpdate        *discordgo.ChannelPinsUpdate        `json:"discord_channel_pins_update,omitempty" dynamodbav:"discord_channel_pins_update,binaryset"`
	ChannelUpdate            *discordgo.ChannelUpdate            `json:"discord_channel_update,omitempty" dynamodbav:"discord_channel_update,binaryset"`
	GuildBanAdd              *discordgo.GuildBanAdd              `json:"discord_guild_ban_add,omitempty" dynamodbav:"discord_guild_ban_add,binaryset"`
	GuildBanRemove           *discordgo.GuildBanRemove           `json:"discord_guild_ban_remove,omitempty" dynamodbav:"discord_guild_ban_remove,binaryset"`
	GuildCreate              *discordgo.GuildCreate              `json:"discord_guild_create,omitempty" dynamodbav:"discord_guild_create,binaryset"`
	GuildDelete              *discordgo.GuildDelete              `json:"discord_guild_delete,omitempty" dynamodbav:"discord_guild_delete,binaryset"`
	GuildEmojisUpdate        *discordgo.GuildEmojisUpdate        `json:"discord_guild_emojis_update,omitempty" dynamodbav:"discord_guild_emojis_update,binaryset"`
	GuildMemberAdd           *discordgo.GuildMemberAdd           `json:"discord_guild_member_add,omitempty" dynamodbav:"discord_guild_member_add,binaryset"`
	GuildMemberRemove        *discordgo.GuildMemberRemove        `json:"discord_guild_member_remove,omitempty" dynamodbav:"discord_guild_member_remove,binaryset"`
	GuildMemberUpdate        *discordgo.GuildMemberUpdate        `json:"discord_guild_member_update,omitempty" dynamodbav:"discord_guild_member_update,binaryset"`
	GuildMembersChunk        *discordgo.GuildMembersChunk        `json:"discord_guild_members_chunk,omitempty" dynamodbav:"discord_guild_members_chunk,binaryset"`
	GuildRoleCreate          *discordgo.GuildRoleCreate          `json:"discord_guild_role_create,omitempty" dynamodbav:"discord_guild_role_create,binaryset"`
	GuildRoleDelete          *discordgo.GuildRoleDelete          `json:"discord_guild_role_delete,omitempty" dynamodbav:"discord_guild_role_delete,binaryset"`
	GuildRoleUpdate          *discordgo.GuildRoleUpdate          `json:"discord_guild_role_update,omitempty" dynamodbav:"discord_guild_role_update,binaryset"`
	GuildUpdate              *discordgo.GuildUpdate              `json:"discord_guild_update,omitempty" dynamodbav:"discord_guild_update,binaryset"`
	MessageCreate            *discordgo.MessageCreate            `json:"discord_message_create,omitempty" dynamodbav:"discord_message_create,binaryset"`
	MessageDelete            *discordgo.MessageDelete            `json:"discord_message_delete,omitempty" dynamodbav:"discord_message_delete,binaryset"`
	MessageReactionAdd       *discordgo.MessageReactionAdd       `json:"discord_message_reaction_add,omitempty" dynamodbav:"discord_message_reaction_add,binaryset"`
	MessageReactionRemove    *discordgo.MessageReactionRemove    `json:"discord_message_reaction_remove,omitempty" dynamodbav:"discord_message_reaction_remove,binaryset"`
	MessageReactionRemoveAll *discordgo.MessageReactionRemoveAll `json:"discord_message_reaction_remove_all,omitempty" dynamodbav:"discord_message_reaction_remove_all,binaryset"`
	MessageUpdate            *discordgo.MessageUpdate            `json:"discord_message_update,omitempty" dynamodbav:"discord_message_update,binaryset"`
	PresenceUpdate           *discordgo.PresenceUpdate           `json:"discord_presence_update,omitempty" dynamodbav:"discord_presence_update,binaryset"`
	GuildIntegrationsUpdate  *discordgo.GuildIntegrationsUpdate  `json:"discord_guild_integrations_update,omitempty" dynamodbav:"discord_guild_integrations_update,binaryset"`
	MessageDeleteBulk        *discordgo.MessageDeleteBulk        `json:"discord_message_delete_bulk,omitempty" dynamodbav:"discord_message_delete_bulk,binaryset"`
	UserUpdate               *discordgo.UserUpdate               `json:"discord_user_update,omitempty" dynamodbav:"discord_user_update,binaryset"`
	VoiceStateUpdate         *discordgo.VoiceStateUpdate         `json:"discord_voice_state_update,omitempty" dynamodbav:"discord_voice_state_update,binaryset"`
	VoiceServerUpdate        *discordgo.VoiceServerUpdate        `json:"discord_voice_server_update,omitempty" dynamodbav:"discord_voice_server_update,binaryset"`
	WebhooksUpdate           *discordgo.WebhooksUpdate           `json:"discord_webhooks_update,omitempty" dynamodbav:"discord_webhooks_update,binaryset"`
	// Cacophony Event Data
	BucketUpdate *BucketUpdate `json:"cacophony_bucket_update,omitempty" dynamodbav:"cacophony_bucket_update,binaryset"`

	// non marshalled events

	ctx            context.Context
	discordSession *discord.Session
	logger         *zap.Logger
	localisations  []interfaces.Localisation
	command        bool
	dm             bool
	prefix         string
	fields         []string
	tokens         map[string]string
	state          *state.State
	botOwnerIDs    []string
	paginator      *paginator.Paginator
	redisClient    *redis.Client
}

// BucketType specifies the scope of a bucket
type BucketType string

// defines various bucket scopes
const (
	GuildBucketType   BucketType = "guild_bucket"
	ChannelBucketType BucketType = "channel_bucket"
	UserBucketType    BucketType = "user_bucket"
)

// BucketUpdate represents Bucket Update event data
type BucketUpdate struct {
	Type      BucketType
	Tag       string
	GuildID   string
	Values    []string
	KeySuffix string
}
