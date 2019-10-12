package events

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

// TODO: fix event struct hashing

// GenerateEventFromDiscordgoEvent generates an Event from a Discordgo Event
func GenerateEventFromDiscordgoEvent(
	botUserID string, eventItem interface{},
) (
	*Event, time.Duration, error,
) {
	event := &Event{
		ReceivedAt: time.Now().UTC(),
		BotUserID:  botUserID,
	}

	expiration := 500 * time.Millisecond

	eventID, err := uuid.NewRandom()
	if err != nil {
		return nil, expiration, err
	}
	event.ID = eventID.String()

	switch t := eventItem.(type) {
	case *discordgo.GuildCreate:
		event.Type = GuildCreateType
		event.GuildCreate = t
		if t.Guild != nil {
			event.GuildID = t.Guild.ID
		}
		event.CacheKey, err = hash(string(event.Type) + guildKey(t.Guild))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildUpdate:
		event.Type = GuildUpdateType
		event.GuildUpdate = t
		if t.Guild != nil {
			event.GuildID = t.Guild.ID
		}
		event.CacheKey, err = hash(string(event.Type) + guildKey(t.Guild))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildDelete:
		event.Type = GuildDeleteType
		event.GuildDelete = t
		if t.Guild != nil {
			event.GuildID = t.Guild.ID
		}
		event.CacheKey, err = hash(string(event.Type) + guildKey(t.Guild))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildMemberAdd:
		event.Type = GuildMemberAddType
		event.GuildMemberAdd = t
		event.GuildID = t.GuildID
		if t.User != nil {
			event.UserID = t.User.ID
		}
		event.CacheKey, err = hash(string(event.Type) + memberKey(t.Member))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildMemberUpdate:
		event.Type = GuildMemberUpdateType
		event.GuildMemberUpdate = t
		event.GuildID = t.GuildID
		if t.User != nil {
			event.UserID = t.User.ID
		}
		event.CacheKey, err = hash(string(event.Type) + memberKey(t.Member))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildMemberRemove:
		event.Type = GuildMemberRemoveType
		event.GuildMemberRemove = t
		event.GuildID = t.GuildID
		if t.User != nil {
			event.UserID = t.User.ID
		}
		event.CacheKey, err = hash(string(event.Type) + memberKey(t.Member))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildMembersChunk:
		event.Type = GuildMembersChunkType
		event.GuildMembersChunk = t
		event.GuildID = t.GuildID
		event.CacheKey, err = hash(string(event.Type) + guildMemberChunkKey(t))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildRoleCreate:
		event.Type = GuildRoleCreateType
		event.GuildRoleCreate = t
		event.GuildID = t.GuildID
		event.CacheKey, err = hash(string(event.Type) + roleKey(t.Role) + t.GuildID)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildRoleUpdate:
		event.Type = GuildRoleUpdateType
		event.GuildRoleUpdate = t
		event.GuildID = t.GuildID
		event.CacheKey, err = hash(string(event.Type) + roleKey(t.Role) + t.GuildID)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildRoleDelete:
		event.Type = GuildRoleDeleteType
		event.GuildRoleDelete = t
		event.GuildID = t.GuildID
		event.CacheKey, err = hash(string(event.Type) + t.RoleID + t.GuildID)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildEmojisUpdate:
		event.Type = GuildEmojisUpdateType
		event.GuildEmojisUpdate = t
		event.GuildID = t.GuildID
		event.CacheKey, err = hash(string(event.Type) + emojisUpdateKey(t))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.ChannelCreate:
		event.Type = ChannelCreateType
		event.ChannelCreate = t
		event.GuildID = t.GuildID
		if t.Channel != nil {
			event.ChannelID = t.Channel.ID
		}
		event.CacheKey, err = hash(string(event.Type) + channelKey(t.Channel))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.ChannelUpdate:
		event.Type = ChannelUpdateType
		event.ChannelUpdate = t
		event.GuildID = t.GuildID
		if t.Channel != nil {
			event.ChannelID = t.Channel.ID
		}
		event.CacheKey, err = hash(string(event.Type) + channelKey(t.Channel))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.ChannelDelete:
		event.Type = ChannelDeleteType
		event.ChannelDelete = t
		event.GuildID = t.GuildID
		if t.Channel != nil {
			event.ChannelID = t.Channel.ID
		}
		event.CacheKey, err = hash(string(event.Type) + channelKey(t.Channel))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.MessageCreate:
		event.Type = MessageCreateType
		event.MessageCreate = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		if t.Author != nil {
			event.UserID = t.Author.ID
		}
		event.MessageID = t.ID
		event.CacheKey, err = hash(string(event.Type) + t.ID)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.MessageUpdate:
		event.Type = MessageUpdateType
		event.MessageUpdate = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		if t.Author != nil {
			event.UserID = t.Author.ID
		}
		event.MessageID = t.ID
		event.CacheKey, err = hash(string(event.Type) + messageUpdateKey(t))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.MessageDelete:
		event.Type = MessageDeleteType
		event.MessageDelete = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		if t.Author != nil {
			event.UserID = t.Author.ID
		}
		event.MessageID = t.ID
		event.CacheKey, err = hash(string(event.Type) + t.ID)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.ChannelPinsUpdate:
		event.Type = ChannelPinsUpdateType
		event.ChannelPinsUpdate = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		event.CacheKey, err = hash(string(event.Type) + t.ChannelID + t.GuildID + t.LastPinTimestamp)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildBanAdd:
		event.Type = GuildBanAddType
		event.GuildBanAdd = t
		event.GuildID = t.GuildID
		if t.User != nil {
			event.UserID = t.User.ID
		}
		event.CacheKey, err = hash(string(event.Type) + t.GuildID + userKey(t.User))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildBanRemove:
		event.Type = GuildBanRemoveType
		event.GuildBanRemove = t
		event.GuildID = t.GuildID
		if t.User != nil {
			event.UserID = t.User.ID
		}
		event.CacheKey, err = hash(string(event.Type) + t.GuildID + userKey(t.User))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.MessageReactionAdd:
		event.Type = MessageReactionAddType
		event.MessageReactionAdd = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		event.UserID = t.UserID
		event.MessageID = t.MessageID
		event.CacheKey, err = hash(string(event.Type) + messageReactionKey(t.MessageReaction))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.MessageReactionRemove:
		event.Type = MessageReactionRemoveType
		event.MessageReactionRemove = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		event.UserID = t.UserID
		event.MessageID = t.MessageID
		event.CacheKey, err = hash(string(event.Type) + messageReactionKey(t.MessageReaction))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.MessageReactionRemoveAll:
		event.Type = MessageReactionRemoveAllType
		event.MessageReactionRemoveAll = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		event.UserID = t.UserID
		event.MessageID = t.MessageID
		event.CacheKey, err = hash(string(event.Type) + messageReactionKey(t.MessageReaction))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.GuildIntegrationsUpdate:
		event.Type = GuildIntegrationsUpdateType
		event.GuildIntegrationsUpdate = t
		event.GuildID = t.GuildID
		event.CacheKey, err = hash(string(event.Type) + t.GuildID)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.PresenceUpdate:
		event.Type = PresenceUpdateType
		event.PresenceUpdate = t
		event.GuildID = t.GuildID
		if t.User != nil {
			event.UserID = t.User.ID
		}
		event.CacheKey, err = hash(string(event.Type) + presenceUpdateKey(t))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.MessageDeleteBulk:
		event.Type = MessageDeleteBulkType
		event.MessageDeleteBulk = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		event.CacheKey, err = hash(string(event.Type) + t.ChannelID + t.GuildID + strings.Join(t.Messages, ""))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.UserUpdate:
		event.Type = UserUpdateType
		event.UserUpdate = t
		if t.User != nil {
			event.UserID = t.User.ID
		}
		event.CacheKey, err = hash(string(event.Type) + userKey(t.User))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.VoiceStateUpdate:
		event.Type = VoiceStateUpdateType
		event.VoiceStateUpdate = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		event.UserID = t.UserID
		event.CacheKey, err = hash(string(event.Type) + voiceStateKey(t.VoiceState))
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.VoiceServerUpdate:
		event.Type = VoiceServerUpdateType
		event.VoiceServerUpdate = t
		event.GuildID = t.GuildID
		event.CacheKey, err = hash(string(event.Type) + t.GuildID + t.Token + t.Endpoint)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.WebhooksUpdate:
		event.Type = WebhooksUpdateType
		event.WebhooksUpdate = t
		event.GuildID = t.GuildID
		event.ChannelID = t.ChannelID
		event.CacheKey, err = hash(string(event.Type) + t.GuildID + t.ChannelID)
		if err != nil {
			return nil, expiration, err
		}
	case *discordgo.TypingStart, *discordgo.Ready, *discordgo.Event, *discordgo.Connect:
		// ignored events
		return nil, expiration, nil
	}

	if event.Type == "" {
		return nil, expiration, errors.New("received unexpected event")
	}

	return event, expiration, nil
}

func hash(data string) (string, error) {
	md5Hasher := md5.New()

	enc := gob.NewEncoder(md5Hasher)
	err := enc.Encode(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(md5Hasher.Sum(nil)), nil
}
