package state

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
)

type State struct {
	client *redis.Client
	db     *gorm.DB
	botIDs []string
}

func NewSate(client *redis.Client, db *gorm.DB, botIDs []string) *State {
	return &State{
		client: client,
		db:     db,
		botIDs: botIDs,
	}
}

// Guild returns the specified Guild from the shard state, returns ErrStateNotFound if not found
func (s *State) Guild(guildID string) (guild *discordgo.Guild, err error) {
	data, err := readStateObject(s.client, guildKey(guildID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &guild)

	if guild != nil {
		membersCount, err := s.client.SCard(guildMembersSetKey(guildID)).Result()
		if err != nil {
			return guild, err
		}
		guild.MemberCount = int(membersCount)

		channelIDs, err := s.client.SMembers(guildChannelsSetKey(guildID)).Result()
		if err != nil {
			return guild, err
		}
		guild.Channels = make([]*discordgo.Channel, 0, len(channelIDs))
		for _, channelID := range channelIDs {
			channel, err := s.Channel(channelID)
			if err != nil {
				return guild, err
			}

			guild.Channels = append(guild.Channels, channel)
		}

		roleIDs, err := s.client.SMembers(guildRolesSetKey(guildID)).Result()
		if err != nil {
			return guild, err
		}
		guild.Roles = make([]*discordgo.Role, 0, len(roleIDs))
		for _, roleID := range roleIDs {
			role, err := s.Role(guildID, roleID)
			if err != nil {
				return guild, err
			}

			guild.Roles = append(guild.Roles, role)
		}

		guild.Members = nil
		guild.Emojis = nil
		guild.VoiceStates = nil
		guild.Presences = nil
	}

	return
}

// Presence returns the specified Presence from the shard state, returns ErrStateNotFound if not found
func (s *State) Presence(guildID, userID string) (presence *discordgo.Presence, err error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, presence := range guild.Presences {
		if presence.User.ID == userID {
			return presence, nil
		}
	}

	return nil, ErrPresenceStateNotFound
}

// Member returns the specified Member from the shard state, returns ErrStateNotFound if not found
func (s *State) Member(guildID, userID string) (member *discordgo.Member, err error) {
	data, err := readStateObject(s.client, memberKey(guildID, userID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &member)
	return
}

// Role returns the specified Role from the shard state, returns ErrStateNotFound if not found
func (s *State) Role(guildID, roleID string) (role *discordgo.Role, err error) {
	data, err := readStateObject(s.client, roleKey(guildID, roleID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &role)
	return
}

// Channel returns the specified Channel from the shard state, returns ErrStateNotFound if not found
func (s *State) Channel(channelID string) (channel *discordgo.Channel, err error) {
	data, err := readStateObject(s.client, channelKey(channelID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &channel)
	return
}

// Emoji returns the specified Emoji from the shard state, returns ErrStateNotFound if not found
func (s *State) Emoji(guildID, emojiID string) (emoji *discordgo.Emoji, err error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, emoji := range guild.Emojis {
		if emoji.ID == emojiID {
			return emoji, nil
		}
	}

	return nil, ErrEmojiStateNotFound
}

// User returns the specified User from the shard state, returns ErrStateNotFound if not found
func (s *State) User(userID string) (user *discordgo.User, err error) {
	data, err := readStateObject(s.client, userKey(userID))
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(data, &user)
	return
}

// AllGuildIDs returns a list of all Guild IDs from the shared state
func (s *State) AllGuildIDs() (guildIDs []string, err error) {
	return readStateSet(s.client, allGuildIDsSetKey())
}

// AllChannelIDs returns a list of all Channel IDs from the shared state
func (s *State) AllChannelIDs() (guildIDs []string, err error) {
	return readStateSet(s.client, allChannelIDsSetKey())
}

// AllUserIDs returns a list of all User IDs from the shared state
func (s *State) AllUserIDs() (userIDs []string, err error) {
	return readStateSet(s.client, allUserIDsSetKey())
}

// GuildUserIDs returns a list of all User IDs in a specific Guild from the shared state
func (s *State) GuildMembers(guildID string) (userIDs []string, err error) {
	return readStateSet(s.client, guildMembersSetKey(guildID))
}

// IsMember true if the User is a member of the specified Guild
func (s *State) IsMember(guildID, userID string) (isMember bool, err error) {
	isMember, err = s.client.SIsMember(guildMembersSetKey(guildID), userID).Result()
	return isMember, err
}

// GuildBannedUserIDs returns the banned user of a server
// only contains items if the bot has the Ban Members or Administrator permission
func (s *State) GuildBannedUserIDs(guildID string) (userIDs []string, err error) {
	return readStateSet(s.client, guildBannedUserIDsSetKey(guildID))
}

// UserChannelPermissions returns the permission of a user in a channel
func (s *State) UserChannelPermissions(userID, channelID string) (apermissions int, err error) {
	var channel *discordgo.Channel
	channel, err = s.Channel(channelID)
	if err != nil {
		return
	}

	var guild *discordgo.Guild
	guild, err = s.Guild(channel.GuildID)
	if err != nil {
		return
	}

	if userID == guild.OwnerID {
		apermissions = discordgo.PermissionAll
		return
	}

	var member *discordgo.Member
	member, err = s.Member(guild.ID, userID)
	if err != nil {
		return
	}

	return memberChannelPermissions(guild, channel, member), nil
}

// UserPermissions returns the permissions of a user in a guild
func (s *State) UserPermissions(userID, guildID string) (apermissions int, err error) {
	var guild *discordgo.Guild
	guild, err = s.Guild(guildID)
	if err != nil {
		return
	}

	if userID == guild.OwnerID {
		apermissions = discordgo.PermissionAll
		return
	}

	var member *discordgo.Member
	member, err = s.Member(guild.ID, userID)
	if err != nil {
		return
	}

	return memberPermissions(guild, member), nil
}

// ChannelMessages returns the messages of a channel
func (s *State) ChannelMessages(channelID string) (messages []discordgo.Message, err error) {
	var messageDatas []string
	messageDatas, err = readStateList(s.client, messagesListKey(channelID))
	if err != nil {
		return nil, err
	}

	if len(messageDatas) == 0 {
		return messages, nil
	}

	var message discordgo.Message
	for _, messageData := range messageDatas {
		err = jsoniter.UnmarshalFromString(messageData, &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// memberChannelPermissions calculates the permissions for a member in a channel
// Source: https://github.com/bwmarrin/discordgo/blob/develop/restapi.go#L503
func memberChannelPermissions(
	guild *discordgo.Guild,
	channel *discordgo.Channel,
	member *discordgo.Member,
) (apermissions int) {
	userID := member.User.ID

	if userID == guild.OwnerID {
		apermissions = discordgo.PermissionAll
		return
	}

	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			apermissions |= role.Permissions
			break
		}
	}

	for _, role := range guild.Roles {
		for _, roleID := range member.Roles {
			if role.ID == roleID {
				apermissions |= role.Permissions
				break
			}
		}
	}

	if apermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		apermissions |= discordgo.PermissionAll
	}

	// Apply @everyone overrides from the channel.
	for _, overwrite := range channel.PermissionOverwrites {
		if guild.ID == overwrite.ID {
			apermissions &= ^overwrite.Deny
			apermissions |= overwrite.Allow
			break
		}
	}

	denies := 0
	allows := 0

	// Member overwrites can override role overrides, so do two passes
	for _, overwrite := range channel.PermissionOverwrites {
		for _, roleID := range member.Roles {
			if overwrite.Type == "role" && roleID == overwrite.ID {
				denies |= overwrite.Deny
				allows |= overwrite.Allow
				break
			}
		}
	}

	apermissions &= ^denies
	apermissions |= allows

	for _, overwrite := range channel.PermissionOverwrites {
		if overwrite.Type == "member" && overwrite.ID == userID {
			apermissions &= ^overwrite.Deny
			apermissions |= overwrite.Allow
			break
		}
	}

	if apermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		apermissions |= discordgo.PermissionAllChannel
	}

	return apermissions
}

// memberPermissions calculates the permissions for a member in a guild
// Source: https://github.com/bwmarrin/discordgo/blob/develop/restapi.go#L503
func memberPermissions(guild *discordgo.Guild, member *discordgo.Member) (apermissions int) {
	userID := member.User.ID

	if userID == guild.OwnerID {
		apermissions = discordgo.PermissionAll
		return
	}

	for _, role := range guild.Roles {
		if role.ID == guild.ID {
			apermissions |= role.Permissions
			break
		}
	}

	for _, role := range guild.Roles {
		for _, roleID := range member.Roles {
			if role.ID == roleID {
				apermissions |= role.Permissions
				break
			}
		}
	}

	if apermissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		apermissions |= discordgo.PermissionAll
	}

	return apermissions
}
