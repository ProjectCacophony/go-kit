package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/regexp"
	"gitlab.com/Cacophony/go-kit/state"
)

var (
	ErrInvalidMessageLink   = errors.New("invalid message link")
	ErrMessageOnWrongServer = errors.New("message on wrong server")
)

func LookupMessageLink(state *state.State, discord *Session, link string) (*discordgo.Message, error) {
	parts := regexp.DiscordMessageLinkRegexp.FindStringSubmatch(link)
	if len(parts) < 5 {
		return nil, ErrInvalidMessageLink
	}

	guildID := parts[2]
	channelID := parts[3]
	messageID := parts[4]

	message, err := FindMessage(state, discord, channelID, messageID)
	if err != nil {
		return nil, err
	}
	if message.GuildID != guildID {
		return nil, ErrMessageOnWrongServer
	}

	return message, nil
}
