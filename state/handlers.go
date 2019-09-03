package state

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var (
	messagesLimit = 10
)

func (s *State) initGuildBans(session *discordgo.Session, guildID string) (err error) {
	//  check if bot is allowed to see bans
	apermissions, err := s.UserPermissions(session.State.User.ID, guildID)
	if err != nil {
		return err
	}
	if apermissions&discordgo.PermissionBanMembers != discordgo.PermissionBanMembers {
		// fmt.Println("resetting bans for", guildID, "because no permissions")
		// reset ban list if not allowed
		err = deleteStateObject(s.client, guildBannedUserIDsSetKey(guildID))
		if err != nil {
			return err
		}
		err = removeFromStateSet(s.client, guildBannedUserIDInitializedGuildIDsSetKey(), guildID)
		return err
	}

	// have we already cached the guild bans for this guild?
	initializedGuildIDs, err := readStateSet(s.client, guildBannedUserIDInitializedGuildIDsSetKey())
	if err != nil {
		return err
	}

	var guildInitialized bool
	for _, initializedGuildID := range initializedGuildIDs {
		if initializedGuildID == guildID {
			guildInitialized = true
		}
	}

	if guildInitialized {
		// fmt.Println("ignoring initializing bans for", guildID, "because already initialized")
		return
	}

	// reset guild bans
	// fmt.Println("resetting bans for", guildID, "because caching new ones")
	err = deleteStateObject(s.client, guildBannedUserIDsSetKey(guildID))
	if err != nil {
		return err
	}

	// cache new guild bans
	bans, err := session.GuildBans(guildID)
	if err != nil {
		return err
	}

	newSet := make([]string, 0)
	for _, ban := range bans {
		newSet = append(newSet, ban.User.ID)
	}
	err = addToStateSet(s.client, guildBannedUserIDsSetKey(guildID), newSet...)
	if err != nil {
		return err
	}
	// fmt.Println("setting bans for", guildID, ":", strings.Join(newSet, ", "))

	err = addToStateSet(s.client, guildBannedUserIDInitializedGuildIDsSetKey(), guildID)
	return err
}

func (s *State) onReady(_ *discordgo.Session, ready *discordgo.Ready) (err error) {
	fmt.Println("running onReady", ready.User.ID)
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
	fmt.Println("running guildAdd", guild.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// carry over previous guild fields if set
	previousGuild, err := s.Guild(guild.ID)
	if err == nil {
		if len(previousGuild.Roles) > len(guild.Roles) {
			guild.Roles = previousGuild.Roles
		}
		if len(previousGuild.Emojis) > len(guild.Emojis) {
			guild.Emojis = previousGuild.Emojis
		}
		if len(previousGuild.Members) > len(guild.Members) {
			guild.Members = previousGuild.Members
		}
		if len(previousGuild.Presences) > len(guild.Presences) {
			guild.Presences = previousGuild.Presences
		}
		if len(previousGuild.Channels) > len(guild.Channels) {
			guild.Channels = previousGuild.Channels
		}
		if len(previousGuild.VoiceStates) > len(guild.VoiceStates) {
			guild.VoiceStates = previousGuild.VoiceStates
		}
	}

	guildMembers := guild.Members

	guild.Presences = nil
	guild.Members = nil
	guild.MemberCount = 0

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
		s.initGuildBans(gS, gGuildID)
		// if err != nil {
		// 	cache.GetLogger().WithField("module", "state").Errorln("error initializing bans for", gGuildID+":", err.Error())
		// }
	}(session, guild.ID)

	// cache guild channels
	for _, channel := range guild.Channels {
		err = s.channelAdd(channel, true)
		if err != nil {
			return err
		}
	}

	// cache guild members and users
	for _, member := range guildMembers {
		err = s.memberAdd(session, member, true, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *State) guildRemove(session *discordgo.Session, guild *discordgo.Guild) (err error) {
	if guild.Unavailable {
		return nil
	}

	fmt.Println("running guildRemove", guild.ID)
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

	// remove guild
	err = deleteStateObject(s.client, guildKey(guild.ID))
	if err != nil {
		return err
	}
	err = removeFromStateSet(s.client, allGuildIDsSetKey(), guild.ID)
	if err != nil {
		return err
	}

	// remove channels
	for _, channel := range guild.Channels {
		err = deleteStateObject(s.client, channelKey(channel.ID))
		if err != nil {
			return err
		}
		err = removeFromStateSet(s.client, allChannelIDsSetKey(), channel.ID)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *State) memberAdd(session *discordgo.Session, member *discordgo.Member, locked, onlyUpdate bool) (err error) {
	fmt.Println("running memberAdd", member.GuildID, member.User.ID)
	if !locked {
		stateLock.Lock()
		defer stateLock.Unlock()
	}

	// TODO: cache guild and member locally?

	// read previous member if exists
	previousMember, err := s.Member(member.GuildID, member.User.ID)
	if err == nil {
		// carry over previous member fields if set
		if member.JoinedAt == "" {
			member.JoinedAt = previousMember.JoinedAt
		}
		if onlyUpdate {
			if len(member.Roles) <= 0 {
				member.Roles = previousMember.Roles
			}
		}
	}

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
		// init guild bans (async) (could be giving or revoking the bot ban permission)
		go func(gS *discordgo.Session, gGuildID string) {
			s.initGuildBans(gS, gGuildID)
			// if err != nil {
			// 	cache.GetLogger().WithField("module", "state").Errorln("error initializing bans for", gGuildID+":", err.Error())
			// }
		}(session, member.GuildID)
	}

	return nil
}

func (s *State) memberRemove(member *discordgo.Member) (err error) {
	fmt.Println("running memberRemove", member.GuildID, member.User.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// remove member and user
	err = deleteStateObject(s.client, memberKey(member.GuildID, member.User.ID))
	if err != nil {
		return err
	}

	// viable?
	allGuildIDs, err := s.AllGuildIDs()
	if err == nil {
		var isMember bool
		for _, guildID := range allGuildIDs {
			if guildID == member.GuildID {
				continue
			}
			isMember, err = s.IsMember(guildID, member.User.ID)
			if err != nil {
				return err
			}
			if isMember {
				break
			}
		}
		if !isMember {
			err = deleteStateObject(s.client, userKey(member.User.ID))
			if err != nil {
				return err
			}
			err = removeFromStateSet(s.client, allUserIDsSetKey(), member.User.ID)
			if err != nil {
				return err
			}
		}
	}

	err = removeFromStateSet(s.client, guildMembersSetKey(member.GuildID), member.User.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *State) roleAdd(session *discordgo.Session, guildID string, role *discordgo.Role) (err error) {
	fmt.Println("running roleAdd", guildID, role.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read role guild
	previousGuild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	// update previous guild
	var updated bool
	for i, previousRole := range previousGuild.Roles {
		if previousRole.ID == role.ID {
			previousGuild.Roles[i] = role
			updated = true
			break
		}
	}
	if !updated {
		previousGuild.Roles = append(previousGuild.Roles, role)
	}

	// cache guild
	err = updateStateObject(s.client, guildKey(previousGuild.ID), previousGuild)
	if err != nil {
		return err
	}

	if role.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator ||
		role.Permissions&discordgo.PermissionBanMembers == discordgo.PermissionBanMembers {
		// init guild bans (async) (could be giving or revoking the bot ban permission)
		go func(gS *discordgo.Session, gGuildID string) {
			s.initGuildBans(gS, gGuildID)
			// if err != nil {
			// 	cache.GetLogger().WithField("module", "state").Errorln("error initializing bans for", gGuildID+":", err.Error())
			// }
		}(session, guildID)
	}

	return nil
}

func (s *State) roleRemove(guildID, roleID string) (err error) {
	fmt.Println("running roleRemove", guildID, roleID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read role guild
	previousGuild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	// remove role
	for i, previousRole := range previousGuild.Roles {
		if previousRole.ID == roleID {
			previousGuild.Roles = append(previousGuild.Roles[:i], previousGuild.Roles[i+1:]...)
			return nil
		}
	}

	// cache guild
	err = updateStateObject(s.client, guildKey(previousGuild.ID), previousGuild)
	return err
}

func (s *State) emojiAdd(guildID string, emoji *discordgo.Emoji) (err error) {
	fmt.Println("running emojiAdd", guildID, emoji.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read emoji guild
	previousGuild, err := s.Guild(guildID)
	if err != nil {
		return err
	}

	// update previous guild
	var updated bool
	for i, previousEmoji := range previousGuild.Emojis {
		if previousEmoji.ID == emoji.ID {
			previousGuild.Emojis[i] = emoji
			updated = true
			break
		}
	}
	if !updated {
		previousGuild.Emojis = append(previousGuild.Emojis, emoji)
	}

	// cache guild
	err = updateStateObject(s.client, guildKey(previousGuild.ID), previousGuild)
	return err
}

func (s *State) emojisAdd(guildID string, emojis []*discordgo.Emoji) (err error) {
	fmt.Println("running emojisAdd", guildID, len(emojis))
	for _, emoji := range emojis {
		err = s.emojiAdd(guildID, emoji)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *State) channelAdd(channel *discordgo.Channel, locked bool) (err error) {
	fmt.Println("running channelAdd", channel.GuildID, channel.ID)
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

		// cache channel
		err = updateStateObject(s.client, channelKey(channel.ID), channel)
		return err
	}

	if channel.Type != discordgo.ChannelTypeDM && channel.Type != discordgo.ChannelTypeGroupDM {
		// read channel guild
		var previousGuild *discordgo.Guild
		previousGuild, err = s.Guild(channel.GuildID)
		if err != nil {
			return err
		}

		// update guild
		previousGuild.Channels = append(previousGuild.Channels, channel)

		// cache guild
		err = updateStateObject(s.client, guildKey(previousGuild.ID), previousGuild)
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
	fmt.Println("running channelRemove", channel.GuildID, channel.ID)
	stateLock.Lock()
	defer stateLock.Unlock()

	// read channel
	previousChannel, err := s.Channel(channel.ID)
	if err != nil {
		return err
	}

	if channel.Type != discordgo.ChannelTypeDM && channel.Type != discordgo.ChannelTypeGroupDM {
		// read channel guild
		var previousGuild *discordgo.Guild
		previousGuild, err = s.Guild(previousChannel.GuildID)
		if err != nil {
			return err
		}

		// update guild
		for i, previousGuildChannel := range previousGuild.Channels {
			if previousGuildChannel.ID == channel.ID {
				previousGuild.Channels = append(previousGuild.Channels[:i], previousGuild.Channels[i+1:]...)
				break
			}
		}

		// cache guild
		err = updateStateObject(s.client, guildKey(previousGuild.ID), previousGuild)
		if err != nil {
			return err
		}
	}

	// cache channel
	err = deleteStateObject(s.client, channelKey(channel.ID))
	if err != nil {
		return err
	}
	err = removeFromStateSet(s.client, allChannelIDsSetKey(), channel.ID)
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
	err = addToStateSet(s.client, guildBannedUserIDsSetKey(guildID), user.ID)
	return err
}

func (s *State) banRemove(guildID string, user *discordgo.User) (err error) {
	err = removeFromStateSet(s.client, guildBannedUserIDsSetKey(guildID), user.ID)
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
		err = s.guildRemove(session, t.Guild)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildDelete guildRemove")
		}
		return nil
	case *discordgo.GuildMemberAdd:
		err = s.memberAdd(session, t.Member, false, false)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildMemberAdd memberAdd")
		}
		return nil
	case *discordgo.GuildMemberUpdate:
		err = s.memberAdd(session, t.Member, false, false)
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
		fmt.Printf("got GuildMembersChunk %s with %d members\n", t.GuildID, len(t.Members))
		for i := range t.Members {
			t.Members[i].GuildID = t.GuildID
			err := s.memberAdd(session, t.Members[i], false, false)
			if err != nil {
				return errors.Wrap(err, "failed to process GuildMembersChunk memberAdd")
			}
		}
		return nil
	case *discordgo.GuildRoleCreate:
		err = s.roleAdd(session, t.GuildID, t.Role)
		if err != nil {
			return errors.Wrap(err, "failed to process GuildRoleCreate roleAdd")
		}
		return nil
	case *discordgo.GuildRoleUpdate:
		err = s.roleAdd(session, t.GuildID, t.Role)
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
		err = s.emojisAdd(t.GuildID, t.Emojis)
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
				Nick:    t.Nick,
				User:    t.User,
				Roles:   t.Roles,
			}

		} else {
			if (t.Nick == "" || t.Nick == previousMember.Nick) &&
				(t.User.Username == "" || t.User.Username == previousMember.User.Username) &&
				(t.User.Discriminator == "" || t.User.Discriminator == previousMember.User.Discriminator) &&
				(t.User.Avatar == "" || t.User.Avatar == previousMember.User.Avatar) &&
				sliceMatches(t.Roles, previousMember.Roles) {
				// fmt.Println("skipped presenceUpdate, no changes")
				return nil
			}

			if t.Nick != "" {
				previousMember.Nick = t.Nick
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

			// PresenceUpdates always contain a list of roles, so there's no need to check for an empty list here
			previousMember.Roles = t.Roles
		}

		err = s.memberAdd(session, previousMember, false, false)
		if err != nil {
			return errors.Wrap(err, "failed to process PresenceUpdate memberAdd")
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

func sliceMatches(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, itemA := range a {
		if !sliceContains(itemA, b) {
			return false
		}
	}

	return true
}

func sliceContains(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}
