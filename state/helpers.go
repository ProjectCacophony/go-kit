package state

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

// nolint: gochecknoglobals
var (
	// MentionRegex matches Discord User Mentions
	// Source: https://github.com/b1naryth1ef/disco/blob/master/disco/bot/command.py#L15, modified to accept IDs
	MentionRegex = regexp.MustCompile("(<@!?)?([0-9]+)(>)?")

	// RoleRegex matches Discord Role Mentions
	// Source: https://github.com/b1naryth1ef/disco/blob/master/disco/bot/command.py#L16, modified to accept IDs
	RoleRegex = regexp.MustCompile("(<@&)?([0-9]+)(>)?")

	// ChannelRegex matches Discord Channel Mentions
	// Source: https://github.com/b1naryth1ef/disco/blob/master/disco/bot/command.py#L17, modified to accept IDs
	ChannelRegex = regexp.MustCompile("(<#)?([0-9]+)(>)?")
)

// UserFromMention finds an user in a mention, can be direct ID input
func (s *State) UserFromMention(mention string) (*discordgo.User, error) {
	result := MentionRegex.FindStringSubmatch(mention)
	if len(result) != 4 {
		return nil, ErrUserNotFound
	}

	return s.User(result[2])
}

// ChannelFromMention finds a text channel on the same server in an mention, can be direct ID input
func (s *State) ChannelFromMention(guildID string, mention string) (*discordgo.Channel, error) {
	result := ChannelRegex.FindStringSubmatch(mention)
	if len(result) != 4 {
		return nil, ErrChannelNotFound
	}

	channel, err := s.Channel(result[2])
	if err != nil {
		return nil, err
	}

	if channel.GuildID != guildID {
		return nil, ErrTargetWrongServer
	}

	if channel.Type != discordgo.ChannelTypeGuildText {
		return nil, ErrTargetWrongType
	}

	return channel, nil
}

// RoleFromMention finds a role in the mention, can be direct ID input
func (s *State) RoleFromMention(guildID string, mention string) (*discordgo.Role, error) {
	result := RoleRegex.FindStringSubmatch(mention)
	if len(result) != 4 {
		return nil, ErrRoleNotFound
	}

	return s.Role(guildID, result[2])
}
