package events

import (
	"context"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/featureflag"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/paginator"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

// Event represents an Event
type Event struct {
	ID         string    `json:"id"`
	CacheKey   string    `json:"cache_key"`
	Type       Type      `json:"type"`
	ReceivedAt time.Time `json:"received_at" dynamodbav:"received_at,unixtime"`
	BotUserID  string    `json:"bot_user_id"`
	GuildID    string    `json:"guild_id,omitempty"`
	ChannelID  string    `json:"channel_id,omitempty"`
	UserID     string    `json:"user_id,omitempty"`
	MessageID  string    `json:"message_id,omitempty"`

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
	BucketUpdate           *BucketUpdate           `json:"cacophony_bucket_update,omitempty"`
	ServerlistServerExpire *ServerlistServerExpire `json:"cacophony_serverlist_server_expire,omitempty"`
	QuickactionRemind      *QuickactionRemind      `json:"cacophony_quickaction_remind,omitempty"`
	QuestionnaireMatch     *QuestionnaireMatch     `json:"cacophony_questionnaire_match,omitempty"`
	AutomodWait            *AutomodWait            `json:"cacophony_automod_wait,omitempty"`

	// non marshalled events
	ctx             context.Context
	discordSession  *discord.Session
	logger          *zap.Logger
	localizations   []interfaces.Localization
	command         bool
	botMention      bool
	dm              bool
	prefix          string
	originalCommand string
	fields          []string
	tokens          map[string]string
	state           *state.State
	paginator       *paginator.Paginator
	redisClient     *redis.Client
	db              *gorm.DB
	httpClient      *http.Client
	questionnaire   *Questionnaire
	storage         *Storage
	featureFlagger  *featureflag.FeatureFlagger
	publisher       *Publisher
}

// BucketType specifies the scope of a bucket
type BucketType string

// defines various bucket scopes
const (
	GuildBucketType   BucketType = "guild_bucket"
	ChannelBucketType BucketType = "channel_bucket"
	UserBucketType    BucketType = "user_bucket"
)

type BucketUpdate struct {
	Type      BucketType
	Tag       string
	GuildID   string
	Values    []string
	KeySuffix string
}

type ServerlistServerExpire struct {
	ID         uint
	GuildID    string
	InviteCode string
	BotID      string
}

type QuickactionRemind struct {
	GuildID   string
	ChannelID string
	MessageID string
	Emoji     *discordgo.Emoji

	ToUserID string
}

type QuestionnaireMatch struct {
	Key     string
	Payload map[string]interface{}
	Filter  QuestionnaireFilter
}

type AutomodWait struct {
	Payload []byte
}
