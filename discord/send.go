package discord

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
)

// SendComplexWithVars sends a message to a given channelID using a given discord session
// it performs various on the actions:
// - pagifies the message content
// - trims the embed
// - escapes @ everyone, and @ here
// TODO: make DMs possible
func SendComplexWithVars(
	session *discordgo.Session,
	localisations []interfaces.Localisation,
	channelID string,
	send *discordgo.MessageSend,
	values ...interface{},
) ([]*discordgo.Message, error) {
	send.Content = localisation.Translate(
		localisations,
		send.Content,
		values...,
	)

	var err error
	var message *discordgo.Message
	var messages []*discordgo.Message
	if send.Embed != nil {
		send.Embed = TrimEmbed(send.Embed)
	}
	send.Content = EscapeDiscordContent(send.Content)
	pages := Pagify(send.Content)
	if len(pages) > 0 {
		for i, page := range pages {
			if i+1 < len(pages) {
				message, err = session.ChannelMessageSend(channelID, page)
			} else {
				send.Content = page
				message, err = session.ChannelMessageSendComplex(channelID, send)
			}
			if err != nil {
				return messages, err
			}
			messages = append(messages, message)
		}

		return messages, nil
	}

	message, err = session.ChannelMessageSendComplex(channelID, send)
	if err != nil {
		return messages, err
	}
	messages = append(messages, message)

	return messages, nil

}
