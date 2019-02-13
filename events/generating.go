package events

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// GenerateRoutingKey generates an Routing Key for AMQP based on a Event Type
func GenerateRoutingKey(eventType Type) string {
	return fmt.Sprintf("cacophony.discord.%s", eventType)
}

// GenerateEventFromDiscordgoEvent generates an Event from a Discordgo Event
// nolint: gocyclo
func GenerateEventFromDiscordgoEvent(botUserID string, eventItem interface{}) (*Event, error) {
	event := &Event{
		ReceivedAt: time.Now().UTC(),
		BotUserID:  botUserID,
	}

	switch t := eventItem.(type) {
	case *discordgo.GuildCreate:
		event.Type = GuildCreateType
		event.GuildCreate = t
	case *discordgo.GuildUpdate:
		event.Type = GuildUpdateType
		event.GuildUpdate = t
	case *discordgo.GuildDelete:
		event.Type = GuildDeleteType
		event.GuildDelete = t
	case *discordgo.GuildMemberAdd:
		event.Type = GuildMemberAddType
		event.GuildMemberAdd = t
	case *discordgo.GuildMemberUpdate:
		event.Type = GuildMemberUpdateType
		event.GuildMemberUpdate = t
	case *discordgo.GuildMemberRemove:
		event.Type = GuildMemberRemoveType
		event.GuildMemberRemove = t
	case *discordgo.GuildMembersChunk:
		event.Type = GuildMembersChunkType
		event.GuildMembersChunk = t
	case *discordgo.GuildRoleCreate:
		event.Type = GuildRoleCreateType
		event.GuildRoleCreate = t
	case *discordgo.GuildRoleUpdate:
		event.Type = GuildRoleUpdateType
		event.GuildRoleUpdate = t
	case *discordgo.GuildRoleDelete:
		event.Type = GuildRoleDeleteType
		event.GuildRoleDelete = t
	case *discordgo.GuildEmojisUpdate:
		event.Type = GuildEmojisUpdateType
		event.GuildEmojisUpdate = t
	case *discordgo.ChannelCreate:
		event.Type = ChannelCreateType
		event.ChannelCreate = t
	case *discordgo.ChannelUpdate:
		event.Type = ChannelUpdateType
		event.ChannelUpdate = t
	case *discordgo.ChannelDelete:
		event.Type = ChannelDeleteType
		event.ChannelDelete = t
	case *discordgo.MessageCreate:
		event.Type = MessageCreateType
		event.MessageCreate = t
	case *discordgo.MessageUpdate:
		event.Type = MessageUpdateType
		event.MessageUpdate = t
	case *discordgo.MessageDelete:
		event.Type = MessageDeleteType
		event.MessageDelete = t
	case *discordgo.ChannelPinsUpdate:
		event.Type = ChannelPinsUpdateType
		event.ChannelPinsUpdate = t
	case *discordgo.GuildBanAdd:
		event.Type = GuildBanAddType
		event.GuildBanAdd = t
	case *discordgo.GuildBanRemove:
		event.Type = GuildBanRemoveType
		event.GuildBanRemove = t
	case *discordgo.MessageReactionAdd:
		event.Type = MessageReactionAddType
		event.MessageReactionAdd = t
	case *discordgo.MessageReactionRemove:
		event.Type = MessageReactionRemoveType
		event.MessageReactionRemove = t
	case *discordgo.MessageReactionRemoveAll:
		event.Type = MessageReactionRemoveAllType
		event.MessageReactionRemoveAll = t
	case *discordgo.GuildIntegrationsUpdate:
		event.Type = GuildIntegrationsUpdateType
		event.GuildIntegrationsUpdate = t
	case *discordgo.PresenceUpdate:
		event.Type = PresenceUpdateType
		event.PresenceUpdate = t
	case *discordgo.MessageDeleteBulk:
		event.Type = MessageDeleteBulkType
		event.MessageDeleteBulk = t
	case *discordgo.UserUpdate:
		event.Type = UserUpdateType
		event.UserUpdate = t
	case *discordgo.VoiceStateUpdate:
		event.Type = VoiceStateUpdateType
		event.VoiceStateUpdate = t
	case *discordgo.VoiceServerUpdate:
		event.Type = VoiceServerUpdateType
		event.VoiceServerUpdate = t
	case *discordgo.WebhooksUpdate:
		event.Type = WebhooksUpdateType
		event.WebhooksUpdate = t
	case *discordgo.TypingStart, *discordgo.Ready, *discordgo.Event, *discordgo.Connect:
		// ignored events
		return nil, nil
	}

	if event.Type != "" {
		return event, nil
	}

	return nil, errors.New("received unexpected event")
}
