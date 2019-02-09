package events

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord"
)

// Event represents an Event
// TODO: event ID or Key
type Event struct {
	Type       Type
	ReceivedAt time.Time
	BotUserID  string

	// discordgo event data
	ChannelCreate            *discordgo.ChannelCreate
	ChannelDelete            *discordgo.ChannelDelete
	ChannelPinsUpdate        *discordgo.ChannelPinsUpdate
	ChannelUpdate            *discordgo.ChannelUpdate
	GuildBanAdd              *discordgo.GuildBanAdd
	GuildBanRemove           *discordgo.GuildBanRemove
	GuildCreate              *discordgo.GuildCreate
	GuildDelete              *discordgo.GuildDelete
	GuildEmojisUpdate        *discordgo.GuildEmojisUpdate
	GuildMemberAdd           *discordgo.GuildMemberAdd
	GuildMemberRemove        *discordgo.GuildMemberRemove
	GuildMemberUpdate        *discordgo.GuildMemberUpdate
	GuildMembersChunk        *discordgo.GuildMembersChunk
	GuildRoleCreate          *discordgo.GuildRoleCreate
	GuildRoleDelete          *discordgo.GuildRoleDelete
	GuildRoleUpdate          *discordgo.GuildRoleUpdate
	GuildUpdate              *discordgo.GuildUpdate
	MessageCreate            *discordgo.MessageCreate
	MessageDelete            *discordgo.MessageDelete
	MessageReactionAdd       *discordgo.MessageReactionAdd
	MessageReactionRemove    *discordgo.MessageReactionRemove
	MessageReactionRemoveAll *discordgo.MessageReactionRemoveAll
	MessageUpdate            *discordgo.MessageUpdate
	PresenceUpdate           *discordgo.PresenceUpdate
	GuildIntegrationsUpdate  *discordgo.GuildIntegrationsUpdate
	MessageDeleteBulk        *discordgo.MessageDeleteBulk
	UserUpdate               *discordgo.UserUpdate
	VoiceStateUpdate         *discordgo.VoiceStateUpdate
	VoiceServerUpdate        *discordgo.VoiceServerUpdate
	WebhooksUpdate           *discordgo.WebhooksUpdate

	// non marshalled events

	ctx            context.Context
	discordSession *discordgo.Session
}

// Context returns the context for the event
func (e *Event) Context() context.Context {
	if e.ctx == nil {
		e.ctx = context.Background()

		return e.ctx
	}

	return e.ctx
}

// WithContext sets the context for the event
func (e *Event) WithContext(ctx context.Context) {
	e.ctx = ctx
}

// Discord gets the Discord API client for the events bot
func (e *Event) Discord() *discordgo.Session {
	if e.discordSession != nil {
		return e.discordSession
	}

	if e.BotUserID == "" {
		panic("could not create discordgo session, no bot user ID set")
	}

	session, err := discord.NewSession(e.BotUserID)
	if err != nil {
		panic(errors.Wrap(err, "could not create discordgo session"))
	}

	e.discordSession = session

	return e.discordSession
}
