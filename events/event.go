package events

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"go.uber.org/zap"
)

// Event represents an Event
// TODO: event ID or Key
type Event struct {
	ID         string
	Type       Type
	ReceivedAt time.Time
	BotUserID  string

	// discordgo event data
	ChannelCreate            *discordgo.ChannelCreate            `json:",omitempty"`
	ChannelDelete            *discordgo.ChannelDelete            `json:",omitempty"`
	ChannelPinsUpdate        *discordgo.ChannelPinsUpdate        `json:",omitempty"`
	ChannelUpdate            *discordgo.ChannelUpdate            `json:",omitempty"`
	GuildBanAdd              *discordgo.GuildBanAdd              `json:",omitempty"`
	GuildBanRemove           *discordgo.GuildBanRemove           `json:",omitempty"`
	GuildCreate              *discordgo.GuildCreate              `json:",omitempty"`
	GuildDelete              *discordgo.GuildDelete              `json:",omitempty"`
	GuildEmojisUpdate        *discordgo.GuildEmojisUpdate        `json:",omitempty"`
	GuildMemberAdd           *discordgo.GuildMemberAdd           `json:",omitempty"`
	GuildMemberRemove        *discordgo.GuildMemberRemove        `json:",omitempty"`
	GuildMemberUpdate        *discordgo.GuildMemberUpdate        `json:",omitempty"`
	GuildMembersChunk        *discordgo.GuildMembersChunk        `json:",omitempty"`
	GuildRoleCreate          *discordgo.GuildRoleCreate          `json:",omitempty"`
	GuildRoleDelete          *discordgo.GuildRoleDelete          `json:",omitempty"`
	GuildRoleUpdate          *discordgo.GuildRoleUpdate          `json:",omitempty"`
	GuildUpdate              *discordgo.GuildUpdate              `json:",omitempty"`
	MessageCreate            *discordgo.MessageCreate            `json:",omitempty"`
	MessageDelete            *discordgo.MessageDelete            `json:",omitempty"`
	MessageReactionAdd       *discordgo.MessageReactionAdd       `json:",omitempty"`
	MessageReactionRemove    *discordgo.MessageReactionRemove    `json:",omitempty"`
	MessageReactionRemoveAll *discordgo.MessageReactionRemoveAll `json:",omitempty"`
	MessageUpdate            *discordgo.MessageUpdate            `json:",omitempty"`
	PresenceUpdate           *discordgo.PresenceUpdate           `json:",omitempty"`
	GuildIntegrationsUpdate  *discordgo.GuildIntegrationsUpdate  `json:",omitempty"`
	MessageDeleteBulk        *discordgo.MessageDeleteBulk        `json:",omitempty"`
	UserUpdate               *discordgo.UserUpdate               `json:",omitempty"`
	VoiceStateUpdate         *discordgo.VoiceStateUpdate         `json:",omitempty"`
	VoiceServerUpdate        *discordgo.VoiceServerUpdate        `json:",omitempty"`
	WebhooksUpdate           *discordgo.WebhooksUpdate           `json:",omitempty"`
	// Cacophony Event Data
	BucketUpdate *BucketUpdate `json:",omitempty"`

	// non marshalled events

	ctx            context.Context
	discordSession *discordgo.Session
	logger         *zap.Logger
	localisations  []interfaces.Localisation
	command        bool
	dm             bool
	prefix         string
	fields         []string
	tokens         map[string]string
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
