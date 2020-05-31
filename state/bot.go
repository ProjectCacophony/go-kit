package state

import (
	"gitlab.com/Cacophony/go-kit/interfaces"
)

// BotForGuild returns a Bot User ID for the given Guild ID
// matches the Discord Permissions if specified
// if possible one should use BotForChannel with the specific permissions
func (s *State) BotForGuild(
	guildID string,
	permissions ...int,
) (
	botID string,
	err error,
) {
	botIDs, err := readStateSet(s.client, guildBotIDsSetKey(guildID))
	if err != nil {
		return "", err
	}

	var permissionsMatch bool
	var botPermissions int

	for _, botID := range botIDs {
		botPermissions, err = s.UserPermissions(
			botID,
			guildID,
		)
		if err != nil {
			continue
		}

		permissionsMatch = true
		for _, permission := range permissions {
			if botPermissions&permission == permission {
				continue
			}

			permissionsMatch = false
		}

		if !permissionsMatch {
			continue
		}

		return botID, nil
	}

	return "", ErrBotForGuildStateNotFound
}

// BotForChannel returns a Bot User ID for the given Channel ID
// matches the Permissions if specified
func (s *State) BotForChannel(
	channelID string,
	permissions ...interfaces.Permission,
) (
	botID string,
	err error,
) {
	channel, err := s.Channel(channelID)
	if err != nil {
		return "", err
	}

	botIDs, err := readStateSet(s.client, guildBotIDsSetKey(channel.GuildID))
	if err != nil {
		return "", err
	}

	var permissionsMatch bool

	for _, botID := range botIDs {
		permissionsMatch = true
		for _, permission := range permissions {
			if permission.Match(s, s.db, botID, channelID, false, false) {
				continue
			}

			permissionsMatch = false
		}

		if !permissionsMatch {
			continue
		}

		return botID, nil
	}

	return "", ErrBotForGuildStateNotFound
}
