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
	ReceivedAt time.Time `json:"received_at"`
	BotUserID  string    `json:"bot_user_id"`
	GuildID    string    `json:"guild_id,omitempty"`
	ChannelID  string    `json:"channel_id,omitempty"`
	UserID     string    `json:"user_id,omitempty"`

	// discordgo event data
	ChannelCreate            *discordgo.ChannelCreate            `json:"discord_channel_create,omitempty"`
	ChannelDelete            *discordgo.ChannelDelete            `json:"discord_channel_delete,omitempty"`
	ChannelPinsUpdate        *discordgo.ChannelPinsUpdate        `json:"discord_channel_pins_update,omitempty"`
	ChannelUpdate            *discordgo.ChannelUpdate            `json:"discord_channel_update,omitempty"`
	GuildBanAdd              *discordgo.GuildBanAdd              `json:"discord_guild_ban_add,omitempty"`
	GuildBanRemove           *discordgo.GuildBanRemove           `json:"discord_guild_ban_remove,omitempty"`
	GuildCreate              *discordgo.GuildCreate              `json:"discord_guild_create,omitempty"`
	GuildDelete              *discordgo.GuildDelete              `json:"discord_guild_delete,omitempty"`
	GuildEmojisUpdate        *discordgo.GuildEmojisUpdate        `json:"discord_guild_emojis_update,omitempty"`
	GuildMemberAdd           *discordgo.GuildMemberAdd           `json:"discord_guild_member_add,omitempty"`
	GuildMemberRemove        *discordgo.GuildMemberRemove        `json:"discord_guild_member_remove,omitempty"`
	GuildMemberUpdate        *discordgo.GuildMemberUpdate        `json:"discord_guild_member_update,omitempty"`
	GuildMembersChunk        *discordgo.GuildMembersChunk        `json:"discord_guild_members_chunk,omitempty"`
	GuildRoleCreate          *discordgo.GuildRoleCreate          `json:"discord_guild_role_create,omitempty"`
	GuildRoleDelete          *discordgo.GuildRoleDelete          `json:"discord_guild_role_delete,omitempty"`
	GuildRoleUpdate          *discordgo.GuildRoleUpdate          `json:"discord_guild_role_update,omitempty"`
	GuildUpdate              *discordgo.GuildUpdate              `json:"discord_guild_update,omitempty"`
	MessageCreate            *discordgo.MessageCreate            `json:"discord_message_create,omitempty"`
	MessageDelete            *discordgo.MessageDelete            `json:"discord_message_delete,omitempty"`
	MessageReactionAdd       *discordgo.MessageReactionAdd       `json:"discord_message_reaction_add,omitempty"`
	MessageReactionRemove    *discordgo.MessageReactionRemove    `json:"discord_message_reaction_remove,omitempty"`
	MessageReactionRemoveAll *discordgo.MessageReactionRemoveAll `json:"discord_message_reaction_remove_all,omitempty"`
	MessageUpdate            *discordgo.MessageUpdate            `json:"discord_message_update,omitempty"`
	PresenceUpdate           *discordgo.PresenceUpdate           `json:"discord_presence_update,omitempty"`
	GuildIntegrationsUpdate  *discordgo.GuildIntegrationsUpdate  `json:"discord_guild_integrations_update,omitempty"`
	MessageDeleteBulk        *discordgo.MessageDeleteBulk        `json:"discord_message_delete_bulk,omitempty"`
	UserUpdate               *discordgo.UserUpdate               `json:"discord_user_update,omitempty"`
	VoiceStateUpdate         *discordgo.VoiceStateUpdate         `json:"discord_voice_state_update,omitempty"`
	VoiceServerUpdate        *discordgo.VoiceServerUpdate        `json:"discord_voice_server_update,omitempty"`
	WebhooksUpdate           *discordgo.WebhooksUpdate           `json:"discord_webhooks_update,omitempty"`
	// Cacophony Event Data
	BucketUpdate *BucketUpdate `json:"cacophony_bucket_update,omitempty"`

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
