package state

import (
	"github.com/bwmarrin/discordgo"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var messagesLimit = 10

func (s *State) onReady(_ *discordgo.Session, ready *discordgo.Ready) (err error) {
	// fmt.Println("running onReady", ready.User.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// cache bot user
	err = updateStateObject(s.client, userKey(ready.User.ID), ready.User)
	if err != nil {
		return err
	}

	// cache private channels
	for _, privateChannel := range ready.PrivateChannels {
		err = updateStateObject(s.client, channelKey(privateChannel.ID), privateChannel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *State) guildAdd(session *discordgo.Session, guild *discordgo.Guild) (err error) {
	// fmt.Println("running guildAdd", guild.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	for _, channel := range guild.Channels {
		err = s.channelAdd(channel, true)
		if err != nil {
			return err
		}
	}

	for _, member := range guild.Members {
		err = s.memberAdd(session, member, true)
		if err != nil {
			return err
		}
	}

	for _, role := range guild.Roles {
		err = s.roleAdd(session, guild.ID, role, true)
		if err != nil {
			return err
		}
	}

	err = s.emojisAdd(guild.ID, guild.Emojis, true)
	if err != nil {
		return err
	}

	guild.Roles = nil
	guild.Emojis = nil
	guild.VoiceStates = nil
	guild.Presences = nil
	guild.Members = nil
	guild.MemberCount = 0
	guild.Channels = nil

	// cache guild
	err = updateStateObject(s.client, guildKey(guild.ID), guild)
	if err != nil {
		return err
	}
	err = addToStateSet(s.client, allGuildIDsSetKey(), guild.ID)
	if err != nil {
		return err
	}
	err = addToStateSet(s.client, guildBotIDsSetKey(guild.ID), session.State.User.ID)
	if err != nil {
		return err
	}

	// init guild bans (async)
	go func(gS *discordgo.Session, gGuildID string) {
		err := s.initGuildBans(gS, gGuildID)
		if err != nil {
			zap.L().Error("error initializing bans", zap.String("guild_id", gGuildID), zap.Error(err))
		}
	}(session, guild.ID)
	// init guild webhooks (async)
	go func(gS *discordgo.Session, gGuildID string) {
		err := s.initGuildWebhooks(gS, gGuildID)
		if err != nil {
			zap.L().Error("error initializing webhooks", zap.String("guild_id", gGuildID), zap.Error(err))
		}
	}(session, guild.ID)
	// init guild invites (async)
	go func(gS *discordgo.Session, gGuildID string) {
		err := s.initGuildInvites(gS, gGuildID)
		if err != nil {
			zap.L().Error("error initializing invites", zap.String("guild_id", gGuildID), zap.Error(err))
		}
	}(session, guild.ID)

	return nil
}

func (s *State) guildRemove(session *discordgo.Session, guild *discordgo.Guild) (err error) {
	if guild.Unavailable {
		return nil
	}

	// fmt.Println("running guildRemove", guild.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// remove this bot for this guild
	err = removeFromStateSet(s.client, guildBotIDsSetKey(guild.ID), session.State.User.ID)
	if err != nil {
		return err
	}

	// check if other bots are on this server, if so, don't do anything further
	botIDs, err := readStateSet(s.client, guildBotIDsSetKey(guild.ID))
	if err != nil {
		return err
	}
	if len(botIDs) > 0 {
		return nil
	}

	err = deleteStateObject(s.client, guildKey(guild.ID))
	if err != nil {
		return err
	}

	err = removeFromStateSet(s.client, allGuildIDsSetKey(), guild.ID)
	if err != nil {
		return err
	}

	err = deleteStateObject(s.client, guildChannelsSetKey(guild.ID))
	if err != nil {
		return err
	}

	err = deleteStateObject(s.client, guildMembersSetKey(guild.ID))
	if err != nil {
		return err
	}

	err = deleteStateObject(s.client, guildEmojiSetKey(guild.ID))
	if err != nil {
		return err
	}

	err = deleteStateObject(s.client, guildRolesSetKey(guild.ID))
	if err != nil {
		return err
	}

	return nil
}

func (s *State) memberAdd(session *discordgo.Session, member *discordgo.Member, locked bool) (err error) {
	// fmt.Println("running memberAdd", member.GuildID, member.User.ID)
	if !locked {
		stateLock.Lock()
		defer stateLock.Unlock()
	}

	// TODO: cache guild and member locally?

	// read previous member if exists
	previousMember, err := s.Member(member.GuildID, member.User.ID)
	if err == nil {
		// carry over previous member fields if set
		if member.JoinedAt.IsZero() {
			member.JoinedAt = previousMember.JoinedAt
		}
	}

	if member.Roles != nil {
		// add new roles
		for _, roleID := range member.Roles {
			err = addToStateSet(s.client, guildMemberRolesSetKey(member.GuildID, member.User.ID), roleID)
			if err != nil {
				return err
			}
		}
		// remove removed roles
		previousRoleIDs, err := s.client.SMembers(guildMemberRolesSetKey(member.GuildID, member.User.ID)).Result()
		if err != nil {
			return err
		}
		for _, previousRoleID := range previousRoleIDs {
			if !sliceContains(previousRoleID, member.Roles) {
				err = removeFromStateSet(s.client, guildMemberRolesSetKey(member.GuildID, member.User.ID), previousRoleID)
				if err != nil {
					return err
				}
			}
		}
	}

	member.Roles = nil

	// cache member and user
	err = updateStateObject(s.client, memberKey(member.GuildID, member.User.ID), member)
	if err != nil {
		return err
	}
	err = updateStateObject(s.client, userKey(member.User.ID), member.User)
	if err != nil {
		return err
	}
	err = addToStateSet(s.client, allUserIDsSetKey(), member.User.ID)
	if err != nil {
		return err
	}

	err = addToStateSet(s.client, guildMembersSetKey(member.GuildID), member.User.ID)
	if err != nil {
		return err
	}

	if member.User.ID == session.State.User.ID {
		// init guild bans (async) (could be giving or revoking the bot bans permission)
		go func(gS *discordgo.Session, gGuildID string) {
			err := s.initGuildBans(gS, gGuildID)
			if err != nil {
				zap.L().Error("error initializing bans", zap.String("guild_id", gGuildID), zap.Error(err))
			}
		}(session, member.GuildID)
		// init guild webhooks (async) (could be giving or revoking the webhooks permission)
		go func(gS *discordgo.Session, gGuildID string) {
			err := s.initGuildWebhooks(gS, gGuildID)
			if err != nil {
				zap.L().Error("error initializing webhooks", zap.String("guild_id", gGuildID), zap.Error(err))
			}
		}(session, member.GuildID)
		// init guild invites (async) (could be giving or revoking the invites permission)
		go func(gS *discordgo.Session, gGuildID string) {
			err := s.initGuildInvites(gS, gGuildID)
			if err != nil {
				zap.L().Error("error initializing invites", zap.String("guild_id", gGuildID), zap.Error(err))
			}
		}(session, member.GuildID)
	}

	return nil
}

func (s *State) memberRemove(member *discordgo.Member) (err error) {
	// fmt.Println("running memberRemove", member.GuildID, member.User.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	err = removeFromStateSet(s.client, guildMembersSetKey(member.GuildID), member.User.ID)
	if err != nil {
		return err
	}

	err = deleteStateObject(s.client, memberKey(member.GuildID, member.User.ID))
	if err != nil {
		return err
	}

	err = deleteStateObject(s.client, guildMemberRolesSetKey(member.GuildID, member.User.ID))
	if err != nil {
		return err
	}

	return nil
}

func (s *State) roleAdd(session *discordgo.Session, guildID string, role *discordgo.Role, locked bool) (err error) {
	// fmt.Println("running roleAdd", guildID, role.ID)
	if !locked {
		stateLock.Lock()
		defer stateLock.Unlock()
	}

	err = updateStateObject(s.client, roleKey(guildID, role.ID), role)
	if err != nil {
		return err
	}

	err = addToStateSet(s.client, guildRolesSetKey(guildID), role.ID)
	if err != nil {
		return err
	}

	// init guild bans (async) (could be giving or revoking the bot bans permission)
	go func(gS *discordgo.Session, gGuildID string) {
		err := s.initGuildBans(gS, gGuildID)
		if err != nil {
			zap.L().Error("error initializing bans", zap.String("guild_id", gGuildID), zap.Error(err))
		}
	}(session, guildID)
	// init guild webhooks (async) (could be giving or revoking the bot webhooks permission)
	go func(gS *discordgo.Session, gGuildID string) {
		err := s.initGuildWebhooks(gS, gGuildID)
		if err != nil {
			zap.L().Error("error initializing webhooks", zap.String("guild_id", gGuildID), zap.Error(err))
		}
	}(session, guildID)
	// init guild invites (async) (could be giving or revoking the bot invites permission)
	go func(gS *discordgo.Session, gGuildID string) {
		err := s.initGuildInvites(gS, gGuildID)
		if err != nil {
			zap.L().Error("error initializing invites", zap.String("guild_id", gGuildID), zap.Error(err))
		}
	}(session, guildID)

	return nil
}

func (s *State) roleRemove(guildID, roleID string) (err error) {
	// fmt.Println("running roleRemove", guildID, roleID)
	stateLock.Lock()
	defer stateLock.Unlock()

	err = removeFromStateSet(s.client, guildRolesSetKey(guildID), roleID)
	if err != nil {
		return err
	}

	err = deleteStateObject(s.client, roleKey(guildID, roleID))
	if err != nil {
		return err
	}

	return nil
}

func (s *State) emojisAdd(guildID string, emojis []*discordgo.Emoji, locked bool) (err error) {
	// fmt.Println("running emojisAdd", guildID, len(emojis))
	if !locked {
		stateLock.Lock()
		defer stateLock.Unlock()
	}

	previousEmojiIDs, err := s.client.SMembers(guildEmojiSetKey(guildID)).Result()
	if err != nil {
		return err
	}

	// add missing emoji
	for _, emoji := range emojis {
		err = updateStateObject(s.client, emojiKey(guildID, emoji.ID), emoji)
		if err != nil {
			return err
		}

		err = addToStateSet(s.client, guildEmojiSetKey(guildID), emoji.ID)
		if err != nil {
			return err
		}
	}

	// remove deleted emoji
	for _, previousEmojiID := range previousEmojiIDs {
		if !emojiContainsEmojiID(previousEmojiID, emojis) {
			err = removeFromStateSet(s.client, guildEmojiSetKey(guildID), previousEmojiID)
			if err != nil {
				return err
			}

			err = deleteStateObject(s.client, emojiKey(guildID, previousEmojiID))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *State) channelAdd(channel *discordgo.Channel, locked bool) (err error) {
	// fmt.Println("running channelAdd", channel.GuildID, channel.ID)
	if !locked {
		stateLock.Lock()
		defer stateLock.Unlock()
	}

	// read channel
	previousChannel, err := s.Channel(channel.ID)
	if err == nil {
		// carry over previous fields if set
		if channel.Messages == nil {
			channel.Messages = previousChannel.Messages
		}
		if channel.PermissionOverwrites == nil {
			channel.PermissionOverwrites = previousChannel.PermissionOverwrites
		}
	}

	if channel.Type != discordgo.ChannelTypeDM && channel.Type != discordgo.ChannelTypeGroupDM {
		err = addToStateSet(s.client, guildChannelsSetKey(channel.GuildID), channel.ID)
		if err != nil {
			return err
		}
	}

	// cache channel
	err = updateStateObject(s.client, channelKey(channel.ID), channel)
	if err != nil {
		return err
	}
	err = addToStateSet(s.client, allChannelIDsSetKey(), channel.ID)
	return err
}

func (s *State) channelRemove(channel *discordgo.Channel) (err error) {
	// fmt.Println("running channelRemove", channel.GuildID, channel.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	if channel.Type != discordgo.ChannelTypeDM && channel.Type != discordgo.ChannelTypeGroupDM {
		err = removeFromStateSet(s.client, guildChannelsSetKey(channel.GuildID), channel.ID)
		if err != nil {
			return err
		}
	}

	err = removeFromStateSet(s.client, allChannelIDsSetKey(), channel.ID)
	if err != nil {
		return err
	}

	err = deleteStateObject(s.client, channelKey(channel.ID))
	return err
}

// func (s *State) presenceAdd(guildID string, presence *discordgo.Presence) (err error) {
// 	fmt.Println("running presenceAdd", guildID, presence.User.ID)
// 	stateLock.Lock()
// 	defer stateLock.Unlock()
//
// 	// read presence guild
// 	previousGuild, err := s.Guild(guildID)
// 	if err != nil {
// 		return err
// 	}
//
// 	// update presence
// 	var updated bool
// 	for i, previousPresence := range previousGuild.Presences {
// 		if previousPresence.User.ID == presence.User.ID {
// 			// Update status
// 			previousGuild.Presences[i].Game = presence.Game
// 			previousGuild.Presences[i].Roles = presence.Roles
// 			if presence.Status != "" {
// 				previousGuild.Presences[i].Status = presence.Status
// 			}
// 			if presence.Nick != "" {
// 				previousGuild.Presences[i].Nick = presence.Nick
// 			}
//
// 			// Update the optionally sent user information
// 			// ID Is a mandatory field so you should not need to check if it is empty
// 			previousGuild.Presences[i].User.ID = presence.User.ID
//
// 			if presence.User.Avatar != "" {
// 				previousGuild.Presences[i].User.Avatar = presence.User.Avatar
// 			}
// 			if presence.User.Discriminator != "" {
// 				previousGuild.Presences[i].User.Discriminator = presence.User.Discriminator
// 			}
// 			if presence.User.Email != "" {
// 				previousGuild.Presences[i].User.Email = presence.User.Email
// 			}
// 			if presence.User.Token != "" {
// 				previousGuild.Presences[i].User.Token = presence.User.Token
// 			}
// 			if presence.User.Username != "" {
// 				previousGuild.Presences[i].User.Username = presence.User.Username
// 			}
//
// 			updated = true
// 		}
// 	}
// 	if !updated {
// 		previousGuild.Presences = append(previousGuild.Presences, presence)
// 	}
//
// 	// cache guild
// 	err = updateStateObject(s.client, guildKey(previousGuild.ID), previousGuild)
// 	return err
// }

func (s *State) banAdd(session *discordgo.Session, guildID string, user *discordgo.User) (err error) {
	// check if bot is allowed to see bans
	apermissions, err := s.UserPermissions(session.State.User.ID, guildID)
	if err != nil {
		return err
	}
	if apermissions&discordgo.PermissionBanMembers != discordgo.PermissionBanMembers {
		return nil
	}

	// add ban
	err = addToStateSet(s.client, guildBanIDsSetKey(guildID), user.ID)
	return err
}

func (s *State) banRemove(guildID string, user *discordgo.User) (err error) {
	err = removeFromStateSet(s.client, guildBanIDsSetKey(guildID), user.ID)
	return err
}

func (s *State) messageCreate(message *discordgo.MessageCreate) (err error) {
	var marshalled string
	marshalled, err = jsoniter.MarshalToString(message)
	if err != nil {
		return err
	}

	err = addToStateList(s.client, messagesListKey(message.ChannelID), marshalled)
	// fmt.Println("added", message.Content, "to", messagesListKey(message.ChannelID))
	if err != nil {
		return err
	}

	err = trimStateList(s.client, messagesListKey(message.ChannelID), int64(messagesLimit-1))
	return err
}

func (s *State) webhookAdd(webhook *discordgo.Webhook) (err error) {
	err = updateStateObject(s.client, webhookKey(webhook.ID), webhook)
	if err != nil {
		return err
	}

	return addToStateSet(s.client, guildWebhookIDsSetKey(webhook.GuildID), webhook.ID)
}

func (s *State) webhookRemove(guildID, webhookID string) (err error) {
	err = deleteStateObject(s.client, webhookKey(webhookID))
	if err != nil {
		return err
	}

	return removeFromStateSet(s.client, guildWebhookIDsSetKey(guildID), webhookID)
}

func (s *State) webhooksUpdate(session *discordgo.Session, guildID, _ string) (err error) {
	stateLock.Lock()
	defer stateLock.Unlock()

	oldWebhooks, err := s.GuildWebhooks(guildID)
	if err != nil {
		return err
	}

	newWebhooks, err := session.GuildWebhooks(guildID)
	if err != nil {
		return err
	}

	for _, oldWebhook := range oldWebhooks {
		if !webhooksSliceContains(oldWebhook.ID, newWebhooks) {
			err = s.webhookRemove(guildID, oldWebhook.ID)
			if err != nil {
				return err
			}
		}
	}

	for _, newWebhooks := range newWebhooks {
		err = s.webhookAdd(newWebhooks)
		if err != nil {
			return err
		}
	}

	return updateStateObject(s.client, guildWebhooksInitializedKey(guildID), true)
}

// SharedStateEventHandler receives events from a discordgo Websocket and updates the shared state with them
func (s *State) SharedStateEventHandler(session *discordgo.Session, i interface{}) error {
	var err error

	switch t := i.(type) {
	case *discordgo.Ready:
		err = s.onReady(session, t)
		if err != nil {
			return errors.Wrap(err, "failed to process OnReady")
		}

		for _, g := range t.Guilds {
			if g.Name == "" {
				continue
			}

			err = s.guildAdd(session, g)
			if err != nil {
				return errors.Wrap(err, "failed to process OnReady guildAdd")
			}
		}

		return nil

	case *discordgo.GuildCreate:
		err = s.guildAdd(session, t.Guild)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildCreate guildAdd")
		}
		return nil
	case *discordgo.GuildUpdate:
		err = s.guildAdd(session, t.Guild)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildUpdate guildAdd")
		}
		return nil
	case *discordgo.GuildDelete:
		if t.Guild.Unavailable || t.Guild.Name == "" {
			return nil
		}

		err = s.guildRemove(session, t.Guild)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildDelete guildRemove")
		}
		return nil
	case *discordgo.GuildMemberAdd:
		err = s.memberAdd(session, t.Member, false)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildMemberAdd memberAdd")
		}
		return nil
	case *discordgo.GuildMemberUpdate:
		err = s.memberAdd(session, t.Member, false)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildMemberUpdate memberAdd")
		}
		return nil
	case *discordgo.GuildMemberRemove:
		err = s.memberRemove(t.Member)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildMemberRemove memberRemove")
		}
		return nil
	case *discordgo.GuildMembersChunk:
		zap.L().Info("received GuildMembersChunk", zap.String("guild_id", t.GuildID), zap.Int("members", len(t.Members)))
		for i := range t.Members {
			t.Members[i].GuildID = t.GuildID
			err := s.memberAdd(session, t.Members[i], false)
			if err != nil {
				return errors.Wrap(err, "failed to process GuildMembersChunk memberAdd")
			}
		}
		return nil
	case *discordgo.GuildRoleCreate:
		err = s.roleAdd(session, t.GuildID, t.Role, false)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildRoleCreate roleAdd")
		}
		return nil
	case *discordgo.GuildRoleUpdate:
		err = s.roleAdd(session, t.GuildID, t.Role, false)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildRoleUpdate roleAdd")
		}
		return nil
	case *discordgo.GuildRoleDelete:
		err = s.roleRemove(t.GuildID, t.RoleID)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildRoleDelete roleRemove")
		}
		return nil
	case *discordgo.GuildEmojisUpdate:
		err = s.emojisAdd(t.GuildID, t.Emojis, false)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildEmojisUpdate emojisAdd")
		}
		return nil
	case *discordgo.ChannelCreate:
		err = s.channelAdd(t.Channel, false)
		if err != nil {
			return errors.Wrap(err, "failed to process ChannelCreate channelAdd")
		}
		return nil
	case *discordgo.ChannelUpdate:
		err = s.channelAdd(t.Channel, false)
		if err != nil {
			return errors.Wrap(err, "failed to process ChannelUpdate channelAdd")
		}
		return nil
	case *discordgo.ChannelDelete:
		err = s.channelRemove(t.Channel)
		if err != nil {
			return errors.Wrap(err, "failed to process ChannelDelete channelRemove")
		}
		return nil
	case *discordgo.GuildBanAdd:
		err = s.banAdd(session, t.GuildID, t.User)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildBanAdd banAdd")
		}
		return nil
	case *discordgo.GuildBanRemove:
		err = s.banRemove(t.GuildID, t.User)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildBanRemove banRemove")
		}
		return nil
	case *discordgo.MessageCreate:
		err = s.messageCreate(t)
		if err != nil {
			return errors.Wrap(err, "failed to process MessageCreate messageCreate")
		}
		return nil
	case *discordgo.PresenceUpdate:
		// err := s.presenceAdd(t.GuildID, &t.Presence)
		// if err != nil {
		// 	return errors.Wrap(err, "failed to process PresenceUpdate presenceAdd")
		// }

		// TODO: store members locally, compare changes?
		previousMember, err := s.Member(t.GuildID, t.User.ID)
		if err != nil {
			// Member not found; this is a user coming online
			previousMember = &discordgo.Member{
				GuildID: t.GuildID,
				User:    t.User,
			}
		} else {
			if (t.User.Username == "" || t.User.Username == previousMember.User.Username) &&
				(t.User.Discriminator == "" || t.User.Discriminator == previousMember.User.Discriminator) &&
				(t.User.Avatar == "" || t.User.Avatar == previousMember.User.Avatar) {
				// fmt.Println("skipped presenceUpdate, no changes")
				return nil
			}

			if t.User.Username != "" {
				previousMember.User.Username = t.User.Username
			}

			if t.User.Discriminator != "" {
				previousMember.User.Discriminator = t.User.Discriminator
			}

			if t.User.Avatar != "" {
				previousMember.User.Avatar = t.User.Avatar
			}
		}

		err = s.memberAdd(session, previousMember, false)
		if err != nil {
			return errors.Wrap(err, "failed to process PresenceUpdate memberAdd")
		}
		return nil
	case *discordgo.WebhooksUpdate:
		err = s.webhooksUpdate(session, t.GuildID, t.ChannelID)
		if err != nil {
			return errors.Wrap(err, "failed to process WebhooksUpdate webhooksUpdate")
		}
		return nil
		/*
		   case *discordgo.MessageUpdate:
		       if s.MaxMessageCount != 0 {
		           err = s.MessageAdd(t.Message)
		       }
		   case *discordgo.MessageDelete:
		       if s.MaxMessageCount != 0 {
		           err = s.MessageRemove(t.Message)
		       }
		   case *discordgo.MessageDeleteBulk:
		       if s.MaxMessageCount != 0 {
		           for _, mID := range t.Messages {
		               s.messageRemoveByID(t.ChannelID, mID)
		           }
		       }
		   case *discordgo.VoiceStateUpdate:
		       if s.TrackVoice {
		           err = s.voiceStateUpdate(t)
		       }
		*/

	}

	return nil
}

func sliceContains(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

func emojiContainsEmojiID(needle string, haystack []*discordgo.Emoji) bool {
	for _, emoji := range haystack {
		if emoji.ID == needle {
			return true
		}
	}

	return false
}
