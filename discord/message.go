package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/state"
)

func FindMessage(
	state *state.State,
	discord *Session,
	channelID string,
	messageID string,
) (
	*discordgo.Message,
	error,
) {
	messages, err := state.ChannelMessages(channelID)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		if messages[i].ID == messageID {
			return &messages[i], nil
		}
	}

	if discord != nil {
		return discord.Client.ChannelMessage(
			channelID,
			messageID,
		)
	}

	return nil, errors.New("message not found")
}
