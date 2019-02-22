package events

import (
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	// nolint: gosec
	"crypto/md5"
)

// TODO: fix event struct hashing

// GenerateRoutingKey generates an Routing Key for AMQP based on a Event Type
func GenerateRoutingKey(eventType Type) string {
	return fmt.Sprintf("cacophony.discord.%s", eventType)
}

// GenerateEventFromDiscordgoEvent generates an Event from a Discordgo Event
// nolint: gocyclo
func GenerateEventFromDiscordgoEvent(botUserID string, eventItem interface{}) (*Event, error) {
	var err error
	event := &Event{
		ReceivedAt: time.Now().UTC(),
		BotUserID:  botUserID,
	}

	switch t := eventItem.(type) {
	case *discordgo.GuildCreate:
		event.Type = GuildCreateType
		event.GuildCreate = t
		event.GuildID = t.Guild.ID
		event.ID, err = hash(t.ID)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildUpdate:
		event.Type = GuildUpdateType
		event.GuildUpdate = t
		event.GuildID = t.Guild.ID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildDelete:
		event.Type = GuildDeleteType
		event.GuildDelete = t
		event.GuildID = t.Guild.ID
		event.ID, err = hash(t.ID)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildMemberAdd:
		event.Type = GuildMemberAddType
		event.GuildMemberAdd = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildMemberUpdate:
		event.Type = GuildMemberUpdateType
		event.GuildMemberUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildMemberRemove:
		event.Type = GuildMemberRemoveType
		event.GuildMemberRemove = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildMembersChunk:
		event.Type = GuildMembersChunkType
		event.GuildMembersChunk = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildRoleCreate:
		event.Type = GuildRoleCreateType
		event.GuildRoleCreate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t.Role.ID)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildRoleUpdate:
		event.Type = GuildRoleUpdateType
		event.GuildRoleUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildRoleDelete:
		event.Type = GuildRoleDeleteType
		event.GuildRoleDelete = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t.RoleID)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildEmojisUpdate:
		event.Type = GuildEmojisUpdateType
		event.GuildEmojisUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.ChannelCreate:
		event.Type = ChannelCreateType
		event.ChannelCreate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t.ID)
		if err != nil {
			return nil, err
		}
	case *discordgo.ChannelUpdate:
		event.Type = ChannelUpdateType
		event.ChannelUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.ChannelDelete:
		event.Type = ChannelDeleteType
		event.ChannelDelete = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t.ID)
		if err != nil {
			return nil, err
		}
	case *discordgo.MessageCreate:
		event.Type = MessageCreateType
		event.MessageCreate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t.ID)
		if err != nil {
			return nil, err
		}
	case *discordgo.MessageUpdate:
		event.Type = MessageUpdateType
		event.MessageUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t.ID + t.Content)
		if err != nil {
			return nil, err
		}
	case *discordgo.MessageDelete:
		event.Type = MessageDeleteType
		event.MessageDelete = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t.ID)
		if err != nil {
			return nil, err
		}
	case *discordgo.ChannelPinsUpdate:
		event.Type = ChannelPinsUpdateType
		event.ChannelPinsUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildBanAdd:
		event.Type = GuildBanAddType
		event.GuildBanAdd = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildBanRemove:
		event.Type = GuildBanRemoveType
		event.GuildBanRemove = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.MessageReactionAdd:
		event.Type = MessageReactionAddType
		event.MessageReactionAdd = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.MessageReactionRemove:
		event.Type = MessageReactionRemoveType
		event.MessageReactionRemove = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.MessageReactionRemoveAll:
		event.Type = MessageReactionRemoveAllType
		event.MessageReactionRemoveAll = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.GuildIntegrationsUpdate:
		event.Type = GuildIntegrationsUpdateType
		event.GuildIntegrationsUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.PresenceUpdate:
		event.Type = PresenceUpdateType
		event.PresenceUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.MessageDeleteBulk:
		event.Type = MessageDeleteBulkType
		event.MessageDeleteBulk = t
		event.GuildID = t.GuildID
		event.ID, err = hash(strings.Join(t.Messages, ""))
		if err != nil {
			return nil, err
		}
	case *discordgo.UserUpdate:
		event.Type = UserUpdateType
		event.UserUpdate = t
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.VoiceStateUpdate:
		event.Type = VoiceStateUpdateType
		event.VoiceStateUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.VoiceServerUpdate:
		event.Type = VoiceServerUpdateType
		event.VoiceServerUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.WebhooksUpdate:
		event.Type = WebhooksUpdateType
		event.WebhooksUpdate = t
		event.GuildID = t.GuildID
		event.ID, err = hash(t)
		if err != nil {
			return nil, err
		}
	case *discordgo.TypingStart, *discordgo.Ready, *discordgo.Event, *discordgo.Connect:
		// ignored events
		return nil, nil
	}

	if event.Type == "" {
		return nil, errors.New("received unexpected event")
	}

	return event, nil
}

func hash(data interface{}) (string, error) {
	// nolint: gosec
	md5Hasher := md5.New()

	enc := gob.NewEncoder(md5Hasher)
	err := enc.Encode(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(md5Hasher.Sum(nil)), nil
}
